package engine_test

import (
	"testing"

	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/schema"
)

func TestDetectMaterialChangeNone(t *testing.T) {
	base := schema.ClassifyRequest{
		Purpose:            "Screen job applicants",
		DataTypes:          []schema.DataType{schema.DataEmployment},
		DeploymentContext:  schema.DeploySaaSB2B,
		AutonomyLevel:      schema.AutonomyDecisionSupport,
		AffectedPopulation: schema.PopJobApplicants,
	}
	current := base
	current.Name = "Renamed only"

	result := engine.DetectMaterialChange(base, current)
	if !result.Material {
		t.Fatal("expected material change for name")
	}
	if result.ConformityImpact != schema.ImpactReassure {
		t.Fatalf("expected reassure, got %s", result.ConformityImpact)
	}
}

func TestDetectMaterialChangeReclassify(t *testing.T) {
	base := schema.ClassifyRequest{
		Purpose:            "Screen job applicants",
		DataTypes:          []schema.DataType{schema.DataEmployment},
		DeploymentContext:  schema.DeploySaaSB2B,
		AutonomyLevel:      schema.AutonomyDecisionSupport,
		AffectedPopulation: schema.PopJobApplicants,
	}
	current := base
	current.AutonomyLevel = schema.AutonomyAutomatedDecision

	result := engine.DetectMaterialChange(base, current)
	if !result.Material {
		t.Fatal("expected material change")
	}
	if result.ConformityImpact != schema.ImpactReclassify {
		t.Fatalf("expected reclassify, got %s", result.ConformityImpact)
	}
}

func TestEUHiringTechnicalControls(t *testing.T) {
	eng, err := engine.Default()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := eng.Classify(schema.ClassifyRequest{
		Purpose:            "Screen job applicants and rank candidates for interview",
		DataTypes:          []schema.DataType{schema.DataPersonal, schema.DataEmployment},
		DeploymentContext:  schema.DeploySaaSB2B,
		AutonomyLevel:      schema.AutonomyDecisionSupport,
		AffectedPopulation: schema.PopJobApplicants,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.TechnicalControls) == 0 {
		t.Fatal("expected EU technical controls on high-risk recruitment")
	}
}
