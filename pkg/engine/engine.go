package engine

import (
	"fmt"
	"strings"

	"github.com/new-world-coder/riskline/pkg/ruleset"
	"github.com/new-world-coder/riskline/pkg/schema"
)

// Engine evaluates a ClassifyRequest against a loaded ruleset.
type Engine struct {
	set *ruleset.Set
}

func New(set *ruleset.Set) *Engine {
	return &Engine{set: set}
}

// Default builds an engine with the embedded EU AI Act ruleset.
func Default() (*Engine, error) {
	set, err := ruleset.LoadDefault()
	if err != nil {
		return nil, err
	}
	return New(set), nil
}

// Classify runs rule matching and returns the highest applicable risk tier.
func (e *Engine) Classify(req schema.ClassifyRequest) (schema.ClassifyResponse, error) {
	if err := validate(req); err != nil {
		return schema.ClassifyResponse{}, err
	}
	if req.GeographicScope == "" {
		req.GeographicScope = schema.GeoEU
	}

	var matched []schema.MatchedRule
	var controls []string
	var judgments []string
	seenControl := map[string]bool{}

	for _, rule := range e.set.Rules {
		if !matches(rule, req) {
			continue
		}
		matched = append(matched, schema.MatchedRule{
			ID:             rule.ID,
			Tier:           schema.RiskTier(rule.Tier),
			ArticleOrAnnex: rule.ArticleOrAnnex,
			Summary:        rule.Summary,
			RulesetVersion: e.set.Version,
			LastUpdated:    e.set.LastUpdated,
		})
		for _, c := range rule.RecommendedControls {
			if !seenControl[c] {
				seenControl[c] = true
				controls = append(controls, c)
			}
		}
		if rule.JudgmentCall && rule.JudgmentCallNote != "" {
			judgments = append(judgments, rule.ID+": "+rule.JudgmentCallNote)
		}
	}

	tier := schema.TierMinimalRisk
	if len(matched) > 0 {
		tier = highestTier(matched)
	} else {
		controls = []string{
			"Document intended purpose and data categories",
			"Re-check classification if the use case expands into Annex III domains",
		}
	}

	return schema.ClassifyResponse{
		Name:                req.Name,
		RiskTier:            tier,
		RulesetVersion:      e.set.Version,
		LastUpdated:         e.set.LastUpdated,
		MatchedRules:        matched,
		Rationale:           buildRationale(tier, matched),
		RecommendedControls: controls,
		JudgmentCalls:       judgments,
		Disclaimer:          schema.Disclaimer,
	}, nil
}

func validate(req schema.ClassifyRequest) error {
	var errs []string
	if strings.TrimSpace(req.Purpose) == "" {
		errs = append(errs, "purpose is required")
	}
	if len(req.DataTypes) == 0 {
		errs = append(errs, "data_types must contain at least one value")
	}
	if req.DeploymentContext == "" {
		errs = append(errs, "deployment_context is required")
	}
	if req.AutonomyLevel == "" {
		errs = append(errs, "autonomy_level is required")
	}
	if req.AffectedPopulation == "" {
		errs = append(errs, "affected_population is required")
	}
	if len(errs) > 0 {
		return &ValidationError{Details: errs}
	}
	return nil
}

// ValidationError lists request problems.
type ValidationError struct {
	Details []string
}

func (e *ValidationError) Error() string {
	return "invalid request: " + strings.Join(e.Details, "; ")
}

func matches(rule ruleset.Rule, req schema.ClassifyRequest) bool {
	w := rule.When
	mode := rule.Match
	if mode == "" {
		mode = "any_group"
	}

	type group struct {
		active bool
		ok     bool
	}
	groups := []group{}

	if len(w.AnyFlags) > 0 {
		g := group{active: true}
		for _, f := range w.AnyFlags {
			if flagSet(req, f) {
				g.ok = true
				break
			}
		}
		groups = append(groups, g)
	}
	if len(w.DeploymentContexts) > 0 {
		g := group{active: true, ok: contains(w.DeploymentContexts, string(req.DeploymentContext))}
		groups = append(groups, g)
	}
	if len(w.AffectedPopulations) > 0 {
		g := group{active: true, ok: contains(w.AffectedPopulations, string(req.AffectedPopulation))}
		groups = append(groups, g)
	}
	if len(w.DataTypes) > 0 {
		g := group{active: true}
		for _, dt := range req.DataTypes {
			if contains(w.DataTypes, string(dt)) {
				g.ok = true
				break
			}
		}
		groups = append(groups, g)
	}
	if len(w.AutonomyLevels) > 0 {
		g := group{active: true, ok: contains(w.AutonomyLevels, string(req.AutonomyLevel))}
		groups = append(groups, g)
	}
	if len(w.PurposeKeywords) > 0 {
		g := group{active: true, ok: purposeHasKeyword(req.Purpose, w.PurposeKeywords)}
		groups = append(groups, g)
	}

	if len(groups) == 0 {
		return false
	}

	if mode == "all_groups" {
		for _, g := range groups {
			if g.active && !g.ok {
				return false
			}
		}
		return true
	}

	// any_group: at least one constrained group matches
	for _, g := range groups {
		if g.active && g.ok {
			return true
		}
	}
	return false
}

func flagSet(req schema.ClassifyRequest, name string) bool {
	switch name {
	case "social_scoring":
		return req.SocialScoring
	case "manipulative_techniques":
		return req.ManipulativeTechniques
	case "exploits_vulnerabilities":
		return req.ExploitsVulnerabilities
	case "real_time_remote_biometric_id":
		return req.RealTimeRemoteBiometricID
	case "emotion_recognition_workplace_or_education":
		return req.EmotionRecognitionWorkplaceOrEducation
	case "uses_biometric_categorisation":
		return req.UsesBiometricCategorisation
	default:
		return false
	}
}

func purposeHasKeyword(purpose string, keywords []string) bool {
	p := strings.ToLower(purpose)
	for _, kw := range keywords {
		if strings.Contains(p, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func highestTier(rules []schema.MatchedRule) schema.RiskTier {
	rank := map[schema.RiskTier]int{
		schema.TierMinimalRisk: 1,
		schema.TierLimitedRisk: 2,
		schema.TierHighRisk:    3,
		schema.TierProhibited:  4,
	}
	best := schema.TierMinimalRisk
	for _, r := range rules {
		if rank[r.Tier] > rank[best] {
			best = r.Tier
		}
	}
	return best
}

func buildRationale(tier schema.RiskTier, matched []schema.MatchedRule) string {
	if len(matched) == 0 {
		return "No prohibited, high-risk, or limited-risk rules matched the described system under the current ruleset. Treated as minimal risk — re-evaluate if the purpose, population, or data types change."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Classified as %s because %d rule(s) matched. Highest-severity matches: ", tier, len(matched))
	parts := make([]string, 0, 3)
	for _, r := range matched {
		if r.Tier != tier {
			continue
		}
		parts = append(parts, r.ArticleOrAnnex+" ("+r.ID+")")
		if len(parts) >= 3 {
			break
		}
	}
	b.WriteString(strings.Join(parts, "; "))
	b.WriteString(".")
	return b.String()
}
