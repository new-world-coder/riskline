package schema

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestRuntimeObservationJSONRoundTrip(t *testing.T) {
	ts := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	orig := RuntimeObservation{
		SystemID:             "hiring-assist-prod",
		DeploymentID:         "deploy-us-east-1",
		ModelID:              "gpt-4.1",
		SystemPromptHash:     "sha256:abc123",
		Tools:                []string{"ats_lookup", "calendar"},
		EventType:            RuntimeEventInference,
		DataCategories:       []DataType{DataPersonal, DataEmployment},
		Geography:            GeoEU,
		AutonomyLevel:        AutonomyDecisionSupport,
		HumanApprovalGranted: true,
		PolicyVersion:        "runtime-policy-v1",
		RulesetVersion:       "eu-ai-act-2024-v0.1.0",
		Timestamp:            ts,
		Metrics: map[string]float64{
			"latency_ms": 42.5,
			"tokens":     128,
		},
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded RuntimeObservation
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("round-trip mismatch:\n  got  %+v\n  want %+v", decoded, orig)
	}
}

func TestRuntimePolicyJSONRoundTrip(t *testing.T) {
	orig := RuntimePolicy{
		ApprovedModels:          []string{"gpt-4.1", "claude-3.5"},
		ForbiddenModels:         []string{"gpt-3.5-turbo"},
		ApprovedTools:             []string{"ats_lookup"},
		ForbiddenTools:            []string{"shell_exec"},
		ApprovedDataCategories:    []DataType{DataPersonal, DataEmployment},
		ForbiddenDataCategories:   []DataType{DataBiometric},
		ApprovedGeography:         []GeographicScope{GeoEU, GeoEUAndGlobal},
		ForbiddenGeography:        []GeographicScope{GeoUnknown},
		MaxAutonomy:               AutonomyDecisionSupport,
		RequiredControls:          []string{"human-review-employment", "audit-log"},
		RiskTier:                  TierHighRisk,
		HumanApprovalRequired:     true,
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded RuntimePolicy
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(orig, decoded) {
		t.Fatalf("round-trip mismatch:\n  got  %+v\n  want %+v", decoded, orig)
	}
}
