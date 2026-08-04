package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/new-world-coder/riskline/pkg/schema"
)

func TestClassifyTable(t *testing.T) {
	eng, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		req  schema.ClassifyRequest
		tier schema.RiskTier
		rule string // expected matched rule id (optional)
	}{
		{
			name: "prohibited social scoring",
			req: schema.ClassifyRequest{
				Purpose:            "Score citizens based on social behaviour",
				DataTypes:          []schema.DataType{schema.DataPersonal},
				DeploymentContext:  schema.DeployPublicAuthority,
				AutonomyLevel:      schema.AutonomyAutomatedDecision,
				AffectedPopulation: schema.PopGeneralPublic,
				SocialScoring:      true,
			},
			tier: schema.TierProhibited,
			rule: "prohibited-social-scoring",
		},
		{
			name: "prohibited manipulation",
			req: schema.ClassifyRequest{
				Purpose:                "Covertly nudge users into purchases",
				DataTypes:              []schema.DataType{schema.DataPersonal},
				DeploymentContext:      schema.DeploySaaSB2C,
				AutonomyLevel:          schema.AutonomyAutonomousAction,
				AffectedPopulation:     schema.PopCustomers,
				ManipulativeTechniques: true,
			},
			tier: schema.TierProhibited,
			rule: "prohibited-manipulation",
		},
		{
			name: "high-risk recruitment",
			req: schema.ClassifyRequest{
				Name:               "Hiring Assist",
				Purpose:            "Screen job applicants and rank candidates for interview",
				DataTypes:          []schema.DataType{schema.DataPersonal, schema.DataEmployment},
				DeploymentContext:  schema.DeploySaaSB2B,
				AutonomyLevel:      schema.AutonomyDecisionSupport,
				AffectedPopulation: schema.PopJobApplicants,
			},
			tier: schema.TierHighRisk,
			rule: "high-risk-recruitment",
		},
		{
			name: "high-risk creditworthiness",
			req: schema.ClassifyRequest{
				Purpose:            "Evaluate creditworthiness and produce a credit score",
				DataTypes:          []schema.DataType{schema.DataFinancialCredit},
				DeploymentContext:  schema.DeploySaaSB2B,
				AutonomyLevel:      schema.AutonomyAutomatedDecision,
				AffectedPopulation: schema.PopCreditConsumers,
			},
			tier: schema.TierHighRisk,
			rule: "high-risk-creditworthiness",
		},
		{
			name: "high-risk biometric categorisation",
			req: schema.ClassifyRequest{
				Purpose:                     "Categorise shoppers by inferred demographic attributes from face images",
				DataTypes:                   []schema.DataType{schema.DataBiometric},
				DeploymentContext:           schema.DeployEmbedded,
				AutonomyLevel:               schema.AutonomyAutomatedDecision,
				AffectedPopulation:          schema.PopCustomers,
				UsesBiometricCategorisation: true,
			},
			tier: schema.TierHighRisk,
			rule: "high-risk-biometric-categorisation",
		},
		{
			name: "high-risk law enforcement risk scoring",
			req: schema.ClassifyRequest{
				Purpose:            "Predict risk of reoffending for bail decisions",
				DataTypes:          []schema.DataType{schema.DataLawEnforcement},
				DeploymentContext:  schema.DeployLawEnforcement,
				AutonomyLevel:      schema.AutonomyDecisionSupport,
				AffectedPopulation: schema.PopSuspectsOrAccused,
			},
			tier: schema.TierHighRisk,
			rule: "high-risk-law-enforcement-risk-scoring",
		},
		{
			name: "high-risk employment monitoring",
			req: schema.ClassifyRequest{
				Purpose:            "Monitor employee productivity and inform performance reviews",
				DataTypes:          []schema.DataType{schema.DataEmployment},
				DeploymentContext:  schema.DeploySaaSB2B,
				AutonomyLevel:      schema.AutonomyDecisionSupport,
				AffectedPopulation: schema.PopEmployees,
			},
			tier: schema.TierHighRisk,
			rule: "high-risk-employment-decisions",
		},
		{
			name: "limited-risk chatbot",
			req: schema.ClassifyRequest{
				Purpose:            "Customer support chatbot for product questions",
				DataTypes:          []schema.DataType{schema.DataPersonal},
				DeploymentContext:  schema.DeploySaaSB2C,
				AutonomyLevel:      schema.AutonomyContentGeneration,
				AffectedPopulation: schema.PopCustomers,
			},
			tier: schema.TierLimitedRisk,
			rule: "limited-risk-chatbot-transparency",
		},
		{
			name: "minimal risk internal docs helper",
			req: schema.ClassifyRequest{
				Purpose:            "Summarise internal engineering design docs",
				DataTypes:          []schema.DataType{schema.DataSyntheticNonPersonal},
				DeploymentContext:  schema.DeployInternalOnly,
				AutonomyLevel:      schema.AutonomyContentGeneration,
				AffectedPopulation: schema.PopEmployees,
			},
			tier: schema.TierMinimalRisk,
		},
		{
			name: "realtime biometric in LE is prohibited over high-risk",
			req: schema.ClassifyRequest{
				Purpose:                   "Real-time facial identification in public spaces",
				DataTypes:                 []schema.DataType{schema.DataBiometric},
				DeploymentContext:         schema.DeployLawEnforcement,
				AutonomyLevel:             schema.AutonomyAutomatedDecision,
				AffectedPopulation:        schema.PopSuspectsOrAccused,
				RealTimeRemoteBiometricID: true,
			},
			tier: schema.TierProhibited,
			rule: "prohibited-realtime-remote-biometric-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eng.Classify(tt.req)
			if err != nil {
				t.Fatalf("Classify: %v", err)
			}
			if got.Disclaimer == "" {
				t.Fatal("disclaimer missing")
			}
			if got.RulesetVersion == "" || got.LastUpdated == "" {
				t.Fatal("ruleset versioning fields missing")
			}
			if got.RiskTier != tt.tier {
				t.Fatalf("tier = %s, want %s (matched=%v)", got.RiskTier, tt.tier, ruleIDs(got))
			}
			if tt.rule != "" && !hasRule(got, tt.rule) {
				t.Fatalf("expected rule %s, got %v", tt.rule, ruleIDs(got))
			}
		})
	}
}

func TestDisclaimerAlwaysPresent(t *testing.T) {
	eng, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := eng.Classify(schema.ClassifyRequest{
		Purpose:            "toy example",
		DataTypes:          []schema.DataType{schema.DataOther},
		DeploymentContext:  schema.DeployOther,
		AutonomyLevel:      schema.AutonomyContentGeneration,
		AffectedPopulation: schema.PopOther,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Disclaimer != schema.Disclaimer {
		t.Fatalf("disclaimer = %q", resp.Disclaimer)
	}
}

func TestValidation(t *testing.T) {
	eng, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	_, err = eng.Classify(schema.ClassifyRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*ValidationError)
	if !ok || len(ve.Details) == 0 {
		t.Fatalf("expected ValidationError, got %T %v", err, err)
	}
}

func TestGoldenRecruitment(t *testing.T) {
	eng, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := eng.Classify(schema.ClassifyRequest{
		Name:               "Hiring Assist",
		Purpose:            "Screen job applicants and rank candidates for interview",
		DataTypes:          []schema.DataType{schema.DataPersonal, schema.DataEmployment},
		DeploymentContext:  schema.DeploySaaSB2B,
		AutonomyLevel:      schema.AutonomyDecisionSupport,
		AffectedPopulation: schema.PopJobApplicants,
		GeographicScope:    schema.GeoEU,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := goldenPath(t, "recruitment_high_risk.json")
	want, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(path, got, 0o644); err != nil {
				t.Fatal(err)
			}
			t.Logf("wrote golden file %s — re-run test", path)
			return
		}
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("golden mismatch for %s\n got:\n%s\nwant:\n%s", path, got, want)
	}
}

func goldenPath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	return filepath.Join(root, "testdata", "golden", name)
}

func hasRule(resp schema.ClassifyResponse, id string) bool {
	for _, r := range resp.MatchedRules {
		if r.ID == id {
			return true
		}
	}
	return false
}

func ruleIDs(resp schema.ClassifyResponse) []string {
	out := make([]string, 0, len(resp.MatchedRules))
	for _, r := range resp.MatchedRules {
		out = append(out, r.ID)
	}
	return out
}
