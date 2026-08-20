package engine

import (
	"fmt"
	"strings"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/schema"
)

// reclassifyFields are ClassifyRequest fields that can change the legal tier.
var reclassifyFields = []struct {
	name string
	eq   func(a, b schema.ClassifyRequest) bool
}{
	{"purpose", func(a, b schema.ClassifyRequest) bool {
		return strings.TrimSpace(strings.ToLower(a.Purpose)) == strings.TrimSpace(strings.ToLower(b.Purpose))
	}},
	{"data_types", func(a, b schema.ClassifyRequest) bool { return sameDataTypes(a, b) }},
	{"deployment_context", func(a, b schema.ClassifyRequest) bool { return a.DeploymentContext == b.DeploymentContext }},
	{"autonomy_level", func(a, b schema.ClassifyRequest) bool { return a.AutonomyLevel == b.AutonomyLevel }},
	{"affected_population", func(a, b schema.ClassifyRequest) bool { return a.AffectedPopulation == b.AffectedPopulation }},
	{"uses_biometric_categorisation", boolEq(func(r schema.ClassifyRequest) bool { return r.UsesBiometricCategorisation })},
	{"real_time_remote_biometric_id", boolEq(func(r schema.ClassifyRequest) bool { return r.RealTimeRemoteBiometricID })},
	{"social_scoring", boolEq(func(r schema.ClassifyRequest) bool { return r.SocialScoring })},
	{"emotion_recognition_workplace_or_education", boolEq(func(r schema.ClassifyRequest) bool { return r.EmotionRecognitionWorkplaceOrEducation })},
	{"manipulative_techniques", boolEq(func(r schema.ClassifyRequest) bool { return r.ManipulativeTechniques })},
	{"exploits_vulnerabilities", boolEq(func(r schema.ClassifyRequest) bool { return r.ExploitsVulnerabilities })},
}

func boolEq(f func(schema.ClassifyRequest) bool) func(a, b schema.ClassifyRequest) bool {
	return func(a, b schema.ClassifyRequest) bool { return f(a) == f(b) }
}

func sameDataTypes(a, b schema.ClassifyRequest) bool {
	if len(a.DataTypes) != len(b.DataTypes) {
		return false
	}
	m := map[schema.DataType]int{}
	for _, d := range a.DataTypes {
		m[d]++
	}
	for _, d := range b.DataTypes {
		if m[d] == 0 {
			return false
		}
		m[d]--
	}
	return true
}

// DetectMaterialChange compares baseline and current system descriptions.
func DetectMaterialChange(baseline, current schema.ClassifyRequest) schema.MaterialChangeResult {
	priorFP := assure.HashClassification(baseline)
	currentFP := assure.HashClassification(current)

	var changed []string
	for _, f := range reclassifyFields {
		if !f.eq(baseline, current) {
			changed = append(changed, f.name)
		}
	}

	impact := schema.ImpactNone
	material := len(changed) > 0
	summary := "No material change detected in system description fields that affect classification."

	if material {
		impact = schema.ImpactReclassify
		summary = fmt.Sprintf("Material change in %d field(s): %s. Re-classify and re-verify controls.",
			len(changed), strings.Join(changed, ", "))
	}

	if !material {
		if baseline.Name != current.Name || baseline.GeographicScope != current.GeographicScope {
			impact = schema.ImpactReassure
			material = true
			summary = "Non-tier metadata changed (name or geographic_scope); re-run assurance probes."
			if baseline.Name != current.Name {
				changed = append(changed, "name")
			}
			if baseline.GeographicScope != current.GeographicScope {
				changed = append(changed, "geographic_scope")
			}
		}
	}

	return schema.MaterialChangeResult{
		Material:           material,
		ConformityImpact:   impact,
		ChangedFields:      changed,
		PriorFingerprint:   priorFP,
		CurrentFingerprint: currentFP,
		Summary:            summary,
		Disclaimer:         schema.Disclaimer,
	}
}
