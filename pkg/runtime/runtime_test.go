package runtime_test

import (
	"testing"
	"time"

	"github.com/new-world-coder/riskline/pkg/evidence"
	"github.com/new-world-coder/riskline/pkg/runtime"
	"github.com/new-world-coder/riskline/pkg/schema"
)

func hiringAssistDescription() schema.ClassifyRequest {
	return schema.ClassifyRequest{
		Name:                "Hiring Assist",
		Purpose:             "Screen job applicants and rank candidates for interview",
		DataTypes:           []schema.DataType{schema.DataPersonal, schema.DataEmployment},
		DeploymentContext:   schema.DeploySaaSB2B,
		AutonomyLevel:       schema.AutonomyDecisionSupport,
		AffectedPopulation:  schema.PopJobApplicants,
		GeographicScope:     schema.GeoEU,
		ModelID:             "gpt-4.1",
		SystemPromptHash:    "sha256:prompt-v1",
		Tools:               []string{"ats_lookup", "calendar"},
		HumanApprovalRequired: true,
	}
}

func hiringAssistClassification() schema.ClassifyResponse {
	return schema.ClassifyResponse{
		Name:           "Hiring Assist",
		RiskTier:       schema.TierHighRisk,
		RulesetVersion: "eu-ai-act-2024-v0.1.0",
		TechnicalControls: []schema.TechnicalControl{
			{ID: "eu-art14-human-oversight-hiring", TechnicalHook: "workflow.human_review_required"},
		},
	}
}

func greenAssure() schema.AssureResponse {
	return schema.AssureResponse{
		ConformityState: schema.ConformityGreen,
		Summary:         "all controls verified",
	}
}

func registerFixture(t *testing.T) (schema.RuntimeBaseline, schema.RuntimePolicy) {
	t.Helper()
	reg, err := runtime.Register(schema.RegisterRuntimeRequest{
		SystemID:          "hiring-assist-prod",
		DeploymentID:      "deploy-eu-1",
		SourceDescription: hiringAssistDescription(),
		Classification:    hiringAssistClassification(),
		Assure:            greenAssure(),
	}, time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return reg.Baseline, reg.Policy
}

func baseObservation() schema.RuntimeObservation {
	return schema.RuntimeObservation{
		SystemID:             "hiring-assist-prod",
		DeploymentID:         "deploy-eu-1",
		ModelID:              "gpt-4.1",
		SystemPromptHash:     "sha256:prompt-v1",
		Tools:                []string{"ats_lookup"},
		EventType:            schema.RuntimeEventInference,
		DataCategories:       []schema.DataType{schema.DataPersonal, schema.DataEmployment},
		Geography:            schema.GeoEU,
		AutonomyLevel:        schema.AutonomyDecisionSupport,
		HumanApprovalGranted: true,
		PolicyVersion:        "runtime-policy-v1",
		Timestamp:            time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC),
	}
}

func TestVerifyGreenWithinEnvelope(t *testing.T) {
	baseline, policy := registerFixture(t)
	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: baseObservation(),
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityGreen {
		t.Fatalf("expected green, got %s violations=%+v", resp.Result.ConformityState, resp.Result.Violations)
	}
	if resp.Receipt.ReceiptHash == "" {
		t.Fatal("expected receipt hash")
	}
}

func TestVerifyForbiddenToolRed(t *testing.T) {
	baseline, policy := registerFixture(t)
	policy.ForbiddenTools = []string{"payments.write"}

	obs := baseObservation()
	obs.EventType = schema.RuntimeEventToolCall
	obs.Tools = []string{"payments.write"}

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityRed {
		t.Fatalf("expected red, got %s", resp.Result.ConformityState)
	}
	if resp.Result.RecommendedAction != "block_and_alert" {
		t.Fatalf("expected block_and_alert, got %s", resp.Result.RecommendedAction)
	}
}

func TestVerifyMissingHumanApprovalRed(t *testing.T) {
	baseline, policy := registerFixture(t)
	obs := baseObservation()
	obs.HumanApprovalGranted = false
	obs.EventType = schema.RuntimeEventToolCall

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityRed {
		t.Fatalf("expected red, got %s", resp.Result.ConformityState)
	}
	found := false
	for _, v := range resp.Result.Violations {
		if v.Code == "missing_human_approval" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected missing_human_approval violation")
	}
}

func TestVerifyModelChangeAmber(t *testing.T) {
	baseline, policy := registerFixture(t)
	obs := baseObservation()
	obs.ModelID = "gpt-4.5"

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityAmber {
		t.Fatalf("expected amber, got %s", resp.Result.ConformityState)
	}
}

func TestVerifyUnapprovedToolRed(t *testing.T) {
	baseline, policy := registerFixture(t)
	obs := baseObservation()
	obs.EventType = schema.RuntimeEventToolCall
	obs.Tools = []string{"shell_exec"}

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityRed {
		t.Fatalf("expected red for unapproved tool, got %s", resp.Result.ConformityState)
	}
}

func TestVerifyProhibitedDataCategory(t *testing.T) {
	baseline, policy := registerFixture(t)
	obs := baseObservation()
	obs.DataCategories = []schema.DataType{schema.DataBiometric}

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityAmber {
		t.Fatalf("expected amber for unapproved data category, got %s", resp.Result.ConformityState)
	}
}

func TestRegisterRequiresGreenAssure(t *testing.T) {
	_, err := runtime.Register(schema.RegisterRuntimeRequest{
		SystemID:          "x",
		DeploymentID:      "y",
		SourceDescription: hiringAssistDescription(),
		Classification:    hiringAssistClassification(),
		Assure: schema.AssureResponse{
			ConformityState: schema.ConformityRed,
		},
	}, time.Now().UTC())
	if err == nil {
		t.Fatal("expected error for non-green assure")
	}
}

func TestVerifyConfigChangeToolDriftAmber(t *testing.T) {
	baseline, policy := registerFixture(t)
	obs := baseObservation()
	obs.EventType = schema.RuntimeEventConfigChange
	obs.Tools = []string{"ats_lookup"}

	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: obs,
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if resp.Result.ConformityState != schema.ConformityAmber {
		t.Fatalf("expected amber for config tool drift, got %s", resp.Result.ConformityState)
	}
}

func TestVerifyReceiptChainIntegration(t *testing.T) {
	baseline, policy := registerFixture(t)

	resp1, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: baseObservation(),
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	resp2, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation:         baseObservation(),
		Policy:              &policy,
		Baseline:            &baseline,
		PreviousReceiptHash: resp1.Receipt.ReceiptHash,
	}, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	if resp2.Receipt.PreviousReceiptHash != resp1.Receipt.ReceiptHash {
		t.Fatalf("expected chain link, got %s", resp2.Receipt.PreviousReceiptHash)
	}

	ok, msg := evidence.VerifyReceiptChain([]schema.VerificationReceipt{resp1.Receipt, resp2.Receipt})
	if !ok {
		t.Fatalf("chain broken: %s", msg)
	}
}

func TestSignReceiptBundleEndToEnd(t *testing.T) {
	baseline, policy := registerFixture(t)
	resp, err := runtime.Verify(schema.VerifyRuntimeRequest{
		Observation: baseObservation(),
		Policy:      &policy,
		Baseline:    &baseline,
	}, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	_, priv, err := evidence.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	bundle, err := evidence.SignReceiptBundle(resp.Receipt, priv)
	if err != nil {
		t.Fatal(err)
	}

	ok, msg := evidence.VerifyReceiptBundle(bundle)
	if !ok {
		t.Fatalf("verify signed receipt: %s", msg)
	}
}

func TestDerivePolicyFromDescription(t *testing.T) {
	desc := hiringAssistDescription()
	policy := runtime.DerivePolicy(desc, hiringAssistClassification())
	if policy.HumanApprovalRequired != true {
		t.Fatal("expected human approval required")
	}
	if len(policy.ApprovedTools) != 2 {
		t.Fatalf("expected 2 approved tools, got %d", len(policy.ApprovedTools))
	}
	if policy.RiskTier != schema.TierHighRisk {
		t.Fatalf("expected high_risk, got %s", policy.RiskTier)
	}
}
