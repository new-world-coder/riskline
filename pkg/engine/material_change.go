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

	// Runtime metadata: changes can invalidate evidence at runtime and
	// require re-classification / re-verification.
	{"model_id", func(a, b schema.ClassifyRequest) bool { return a.ModelID == b.ModelID }},
	{"system_prompt_hash", func(a, b schema.ClassifyRequest) bool { return a.SystemPromptHash == b.SystemPromptHash }},
	{"tools", func(a, b schema.ClassifyRequest) bool { return sameStringSet(a.Tools, b.Tools) }},
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

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := map[string]int{}
	for _, d := range a {
		m[d]++
	}
	for _, d := range b {
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
		var nonTierChanged []string
		if baseline.Name != current.Name {
			nonTierChanged = append(nonTierChanged, "name")
		}
		if baseline.GeographicScope != current.GeographicScope {
			nonTierChanged = append(nonTierChanged, "geographic_scope")
		}
		// Human-in-the-loop toggles may change assurance verification,
		// but do not necessarily shift the legal tier.
		if baseline.HumanApprovalRequired != current.HumanApprovalRequired {
			nonTierChanged = append(nonTierChanged, "human_approval_required")
		}

		if len(nonTierChanged) > 0 {
			impact = schema.ImpactReassure
			material = true
			changed = append(changed, nonTierChanged...)
			summary = fmt.Sprintf("Non-tier metadata changed in %d field(s) (%s); re-run assurance probes.",
				len(nonTierChanged), strings.Join(nonTierChanged, ", "))
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
