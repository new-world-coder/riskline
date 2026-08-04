package ruleset

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/eu_ai_act_v0.1.0.json
var euAIActV010 []byte

// RuleWhen describes match conditions. Empty slices mean "don't constrain on this field".
type RuleWhen struct {
	AnyFlags             []string `json:"any_flags"`
	DeploymentContexts   []string `json:"deployment_contexts"`
	AffectedPopulations  []string `json:"affected_populations"`
	DataTypes            []string `json:"data_types"`
	AutonomyLevels       []string `json:"autonomy_levels"`
	PurposeKeywords      []string `json:"purpose_keywords"`
}

// Rule is one versioned classification mapping.
type Rule struct {
	ID                   string   `json:"id"`
	Tier                 string   `json:"tier"`
	ArticleOrAnnex       string   `json:"article_or_annex"`
	Summary              string   `json:"summary"`
	When                 RuleWhen `json:"when"`
	Match                string   `json:"match"` // "any_group" (default) or "all_groups"
	RecommendedControls  []string `json:"recommended_controls"`
	JudgmentCall         bool     `json:"judgment_call"`
	JudgmentCallNote     string   `json:"judgment_call_note"`
}

// Set is a loaded ruleset document.
type Set struct {
	Version     string   `json:"version"`
	LastUpdated string   `json:"last_updated"`
	Framework   string   `json:"framework"`
	Notes       []string `json:"notes"`
	Rules       []Rule   `json:"rules"`
}

// LoadDefault loads the embedded EU AI Act v0.1.0 ruleset.
func LoadDefault() (*Set, error) {
	return Parse(euAIActV010)
}

// Parse decodes a ruleset from JSON bytes.
func Parse(data []byte) (*Set, error) {
	var s Set
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("ruleset: parse: %w", err)
	}
	if s.Version == "" || s.LastUpdated == "" {
		return nil, fmt.Errorf("ruleset: version and last_updated are required")
	}
	if len(s.Rules) == 0 {
		return nil, fmt.Errorf("ruleset: no rules")
	}
	for i, r := range s.Rules {
		if r.ID == "" || r.Tier == "" || r.ArticleOrAnnex == "" {
			return nil, fmt.Errorf("ruleset: rule[%d] missing id/tier/article", i)
		}
	}
	return &s, nil
}
