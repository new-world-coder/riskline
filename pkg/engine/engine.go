package engine

import (
	"fmt"
	"strings"

	"github.com/new-world-coder/riskline/pkg/ruleset"
	"github.com/new-world-coder/riskline/pkg/schema"
)

// Engine evaluates a ClassifyRequest against one or more regime packs.
type Engine struct {
	loader         *ruleset.Loader
	defaultRegimes []string
	// set is retained for New(set) single-pack callers (tests / embedders).
	set *ruleset.Set
}

// New builds a single-pack engine (EU AI Act when set comes from LoadDefault).
func New(set *ruleset.Set) *Engine {
	return &Engine{set: set, defaultRegimes: []string{ruleset.RegimeEUAIAct}}
}

// NewWithLoader builds a multi-pack engine. defaultRegimes is used when the
// request omits regimes (after CLI/API config resolution callers may still
// pass them on the request).
func NewWithLoader(loader *ruleset.Loader, defaultRegimes []string) (*Engine, error) {
	if loader == nil {
		return nil, fmt.Errorf("engine: nil loader")
	}
	regs := ruleset.NormalizeRegimes(defaultRegimes)
	if len(regs) == 0 {
		regs = []string{ruleset.RegimeEUAIAct}
	}
	for _, r := range regs {
		if _, err := loader.Load(r); err != nil {
			return nil, err
		}
	}
	return &Engine{loader: loader, defaultRegimes: regs}, nil
}

// Default builds an engine with the embedded EU AI Act ruleset only.
func Default() (*Engine, error) {
	loader, err := ruleset.DefaultLoader()
	if err != nil {
		return nil, err
	}
	return NewWithLoader(loader, []string{ruleset.RegimeEUAIAct})
}

// Classify runs rule matching for the resolved regimes.
//
// Backward compatible shape: the top-level risk_tier / matched_rules fields
// always reflect the first resolved regime. classifications[] is only filled
// when more than one regime was evaluated, so the historical EU-only JSON
// golden remains byte-identical.
func (e *Engine) Classify(req schema.ClassifyRequest) (schema.ClassifyResponse, error) {
	if err := validate(req); err != nil {
		return schema.ClassifyResponse{}, err
	}
	if req.GeographicScope == "" {
		req.GeographicScope = schema.GeoEU
	}

	regimes, err := e.resolveRegimes(req.Regimes)
	if err != nil {
		return schema.ClassifyResponse{}, err
	}

	results := make([]schema.RegimeClassification, 0, len(regimes))
	for _, regime := range regimes {
		pack, set, character, err := e.packFor(regime)
		if err != nil {
			return schema.ClassifyResponse{}, err
		}
		_ = pack
		rc := classifyAgainst(set, character, regime, req)
		results = append(results, rc)
	}

	primary := results[0]
	resp := schema.ClassifyResponse{
		Name:                req.Name,
		RiskTier:            primary.RiskTier,
		RulesetVersion:      primary.RulesetVersion,
		LastUpdated:         primary.LastUpdated,
		MatchedRules:        primary.MatchedRules,
		Rationale:           primary.Rationale,
		RecommendedControls: primary.RecommendedControls,
		TechnicalControls:   primary.TechnicalControls,
		JudgmentCalls:       primary.JudgmentCalls,
		MappingOnly:         primary.MappingOnly,
		Disclaimer:          schema.Disclaimer,
	}

	// Keep single default EU responses free of additive fields so
	// testdata/golden/recruitment_high_risk.json stays stable.
	if !isSingleDefaultEU(regimes) {
		resp.Regime = primary.Regime
		resp.Classifications = results
	}

	return resp, nil
}

func (e *Engine) resolveRegimes(requested []string) ([]string, error) {
	regs := ruleset.NormalizeRegimes(requested)
	if len(regs) == 0 {
		regs = append([]string(nil), e.defaultRegimes...)
	}
	if len(regs) == 0 {
		regs = []string{ruleset.RegimeEUAIAct}
	}
	for _, r := range regs {
		if _, _, _, err := e.packFor(r); err != nil {
			return nil, &ValidationError{Details: []string{err.Error()}}
		}
	}
	return regs, nil
}

func (e *Engine) packFor(regime string) (*ruleset.Pack, *ruleset.Set, string, error) {
	if e.loader != nil {
		p, err := e.loader.Load(regime)
		if err != nil {
			return nil, nil, "", err
		}
		return p, p.Set, p.Character, nil
	}
	// Legacy New(set) path: only eu-ai-act is available.
	id := ruleset.NormalizeRegime(regime)
	if id != ruleset.RegimeEUAIAct {
		return nil, nil, "", fmt.Errorf("ruleset: unknown regime %q (available: %s)", id, ruleset.RegimeEUAIAct)
	}
	if e.set == nil {
		return nil, nil, "", fmt.Errorf("engine: no ruleset loaded")
	}
	return &ruleset.Pack{ID: id, Character: ruleset.CharacterHardLaw, Set: e.set}, e.set, ruleset.CharacterHardLaw, nil
}

func isSingleDefaultEU(regimes []string) bool {
	return len(regimes) == 1 && regimes[0] == ruleset.RegimeEUAIAct
}

func classifyAgainst(set *ruleset.Set, character, regime string, req schema.ClassifyRequest) schema.RegimeClassification {
	var matched []schema.MatchedRule
	var controls []string
	var techControls []schema.TechnicalControl
	var judgments []string
	seenControl := map[string]bool{}
	seenTech := map[string]bool{}

	for _, rule := range set.Rules {
		if !matches(rule, req) {
			continue
		}
		matched = append(matched, schema.MatchedRule{
			ID:             rule.ID,
			Tier:           schema.RiskTier(rule.Tier),
			ArticleOrAnnex: rule.ArticleOrAnnex,
			Summary:        rule.Summary,
			RulesetVersion: set.Version,
			LastUpdated:    set.LastUpdated,
		})
		for _, c := range rule.RecommendedControls {
			if !seenControl[c] {
				seenControl[c] = true
				controls = append(controls, c)
			}
		}
		for _, tc := range rule.TechnicalControls {
			if tc.ID == "" || seenTech[tc.ID] {
				continue
			}
			seenTech[tc.ID] = true
			techControls = append(techControls, schema.TechnicalControl{
				ID:            tc.ID,
				PaperRef:      tc.PaperRef,
				Summary:       tc.Summary,
				TechnicalHook: tc.TechnicalHook,
				EvidenceType:  tc.EvidenceType,
			})
		}
		if rule.JudgmentCall && rule.JudgmentCallNote != "" {
			judgments = append(judgments, rule.ID+": "+rule.JudgmentCallNote)
		}
	}

	mappingOnly := character == ruleset.CharacterMapping
	tier := schema.TierMinimalRisk
	rationale := buildRationale(tier, matched)

	if mappingOnly {
		rationale = buildMappingRationale(matched)
		if len(controls) == 0 && len(matched) == 0 {
			controls = []string{
				"Document intended purpose and re-run mapping when the use case expands",
			}
		}
	} else if len(matched) > 0 {
		tier = highestTier(matched)
		rationale = buildRationale(tier, matched)
	} else {
		controls = []string{
			"Document intended purpose and data categories",
			"Re-check classification if the use case expands into Annex III domains",
		}
	}

	return schema.RegimeClassification{
		Regime:              regime,
		Character:           character,
		RiskTier:            tier,
		RulesetVersion:      set.Version,
		LastUpdated:         set.LastUpdated,
		MatchedRules:        matched,
		Rationale:           rationale,
		RecommendedControls: controls,
		TechnicalControls:   techControls,
		JudgmentCalls:       judgments,
		MappingOnly:         mappingOnly,
	}
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

func buildMappingRationale(matched []schema.MatchedRule) string {
	if len(matched) == 0 {
		return "No NIST AI RMF mapping rules matched. This is an advisory mapping only — not a US legal risk tier. Re-run when purpose, population, autonomy, or data types change."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "NIST AI RMF mapping: %d subcategory mapping(s) matched. Examples: ", len(matched))
	parts := make([]string, 0, 3)
	for _, r := range matched {
		parts = append(parts, r.ArticleOrAnnex+" ("+r.ID+")")
		if len(parts) >= 3 {
			break
		}
	}
	b.WriteString(strings.Join(parts, "; "))
	b.WriteString(". This is advisory mapping only — not a US legal risk tier.")
	return b.String()
}
