package assure_test

import (
	"testing"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/schema"
)

func TestEvaluateAllPassed(t *testing.T) {
	classification := schema.ClassifyResponse{
		Name:           "Hiring Assist",
		RulesetVersion: "eu-ai-act-2024-v0.1.0",
		TechnicalControls: []schema.TechnicalControl{
			{
				ID:            "eu-art14-human-oversight-hiring",
				TechnicalHook: "workflow.human_review_required",
				PaperRef:      "EU AI Act Article 14",
				EvidenceType:  "audit_log",
			},
		},
	}

	resp := assure.Evaluate(schema.AssureRequest{
		Classification: classification,
		Probes: map[string]bool{
			"workflow.human_review_required": true,
		},
	})

	if resp.ConformityState != schema.ConformityGreen {
		t.Fatalf("expected green, got %s", resp.ConformityState)
	}
	if len(resp.EvidenceRecords) != 1 {
		t.Fatalf("expected 1 evidence record, got %d", len(resp.EvidenceRecords))
	}
	if resp.EvidenceRecords[0].Algorithm != "sha256" {
		t.Fatalf("expected sha256, got %s", resp.EvidenceRecords[0].Algorithm)
	}
}

func TestEvaluateFailedProbe(t *testing.T) {
	classification := schema.ClassifyResponse{
		RulesetVersion: "eu-ai-act-2024-v0.1.0",
		TechnicalControls: []schema.TechnicalControl{
			{ID: "c1", TechnicalHook: "workflow.human_review_required"},
		},
	}

	resp := assure.Evaluate(schema.AssureRequest{
		Classification: classification,
		Probes:         map[string]bool{"workflow.human_review_required": false},
	})

	if resp.ConformityState != schema.ConformityRed {
		t.Fatalf("expected red, got %s", resp.ConformityState)
	}
}

func TestEvidenceChain(t *testing.T) {
	classification := schema.ClassifyResponse{
		RulesetVersion: "eu-ai-act-2024-v0.1.0",
		TechnicalControls: []schema.TechnicalControl{
			{ID: "c1", TechnicalHook: "hook.a", PaperRef: "Art 9", EvidenceType: "log"},
			{ID: "c2", TechnicalHook: "hook.b", PaperRef: "Art 14", EvidenceType: "log"},
		},
	}

	resp := assure.Evaluate(schema.AssureRequest{
		Classification: classification,
		Probes: map[string]bool{
			"hook.a": true,
			"hook.b": true,
		},
	})

	ok, msg := assure.VerifyEvidenceChain(resp.EvidenceRecords)
	if !ok {
		t.Fatalf("chain broken: %s", msg)
	}
	if resp.EvidenceRecords[1].PreviousHash != resp.EvidenceRecords[0].ContentHash {
		t.Fatal("second record should link to first content hash")
	}
}
