package ruleset

import (
	"testing"
)

func TestLoadDefault(t *testing.T) {
	s, err := LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault: %v", err)
	}
	if s.Version != "eu-ai-act-2024-v0.1.0" {
		t.Fatalf("version = %q", s.Version)
	}
	if s.LastUpdated == "" {
		t.Fatal("last_updated empty")
	}
	if len(s.Rules) < 8 {
		t.Fatalf("expected at least 8 rules, got %d", len(s.Rules))
	}

	tiers := map[string]bool{}
	for _, r := range s.Rules {
		tiers[r.Tier] = true
		if r.Summary == "" {
			t.Errorf("rule %s has empty summary", r.ID)
		}
	}
	for _, want := range []string{"prohibited", "high_risk", "limited_risk"} {
		if !tiers[want] {
			t.Errorf("missing tier %s in ruleset", want)
		}
	}
}

func TestParseRejectsUnversioned(t *testing.T) {
	_, err := Parse([]byte(`{"rules":[{"id":"x","tier":"minimal_risk","article_or_annex":"n/a","summary":"x"}]}`))
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestParseRejectsEmptyRules(t *testing.T) {
	_, err := Parse([]byte(`{"version":"v","last_updated":"2026-01-01","rules":[]}`))
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestParseRejectsIncompleteRule(t *testing.T) {
	_, err := Parse([]byte(`{"version":"v","last_updated":"2026-01-01","rules":[{"id":"","tier":"high_risk","article_or_annex":"x","summary":"s"}]}`))
	if err == nil {
		t.Fatal("expected error for incomplete rule")
	}
}

func TestParseRejectsInvalidJSON(t *testing.T) {
	_, err := Parse([]byte(`{`))
	if err == nil {
		t.Fatal("expected parse error")
	}
}
