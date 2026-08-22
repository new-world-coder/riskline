package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/schema"
)

const hashAlgorithm = "sha256"

var autonomyRank = map[schema.AutonomyLevel]int{
	schema.AutonomyContentGeneration: 1,
	schema.AutonomyDecisionSupport:   2,
	schema.AutonomyAutomatedDecision: 3,
	schema.AutonomyAutonomousAction:  4,
}

// Register derives an approved runtime baseline and policy envelope.
func Register(req schema.RegisterRuntimeRequest, now time.Time) (schema.RegisterRuntimeResponse, error) {
	if strings.TrimSpace(req.SystemID) == "" {
		return schema.RegisterRuntimeResponse{}, fmt.Errorf("system_id is required")
	}
	if strings.TrimSpace(req.DeploymentID) == "" {
		return schema.RegisterRuntimeResponse{}, fmt.Errorf("deployment_id is required")
	}
	if req.Assure.ConformityState != schema.ConformityGreen {
		return schema.RegisterRuntimeResponse{}, fmt.Errorf("assure must be green before registration (got %s)", req.Assure.ConformityState)
	}

	desc := req.SourceDescription
	if desc.Purpose == "" {
		return schema.RegisterRuntimeResponse{}, fmt.Errorf("source_description is required")
	}

	if now.IsZero() {
		now = time.Now().UTC()
	}

	policy := DerivePolicy(desc, req.Classification)
	baseline := schema.RuntimeBaseline{
		SystemID:             req.SystemID,
		DeploymentID:         req.DeploymentID,
		Fingerprint:          assure.HashClassification(desc),
		ApprovedAt:           now,
		SourceClassification: desc,
	}

	return schema.RegisterRuntimeResponse{
		Baseline:   baseline,
		Policy:     policy,
		Disclaimer: schema.Disclaimer,
	}, nil
}

// DerivePolicy builds a runtime envelope from the approved system description.
func DerivePolicy(desc schema.ClassifyRequest, classification schema.ClassifyResponse) schema.RuntimePolicy {
	policy := schema.RuntimePolicy{
		MaxAutonomy:           desc.AutonomyLevel,
		RiskTier:              classification.RiskTier,
		HumanApprovalRequired: desc.HumanApprovalRequired,
		ApprovedDataCategories: append([]schema.DataType(nil), desc.DataTypes...),
		RequiredControls:      controlIDs(classification.TechnicalControls),
	}

	if desc.ModelID != "" {
		policy.ApprovedModels = []string{desc.ModelID}
	}
	if len(desc.Tools) > 0 {
		policy.ApprovedTools = append([]string(nil), desc.Tools...)
	}
	if desc.GeographicScope != "" {
		policy.ApprovedGeography = []schema.GeographicScope{desc.GeographicScope}
	}

	return policy
}

func controlIDs(controls []schema.TechnicalControl) []string {
	out := make([]string, 0, len(controls))
	for _, tc := range controls {
		out = append(out, tc.ID)
	}
	sort.Strings(out)
	return out
}

// Verify compares a runtime observation against policy and optional baseline.
func Verify(req schema.VerifyRuntimeRequest, now time.Time) (schema.VerifyRuntimeResponse, error) {
	if req.Policy == nil {
		return schema.VerifyRuntimeResponse{}, fmt.Errorf("policy is required (inline or from registration)")
	}
	if strings.TrimSpace(req.Observation.SystemID) == "" {
		return schema.VerifyRuntimeResponse{}, fmt.Errorf("observation.system_id is required")
	}

	if now.IsZero() {
		now = time.Now().UTC()
	}

	policy := *req.Policy
	obs := req.Observation
	violations := evaluateEnvelope(obs, policy)

	if req.Baseline != nil {
		current := observationAsClassifyRequest(obs, req.Baseline.SourceClassification)
		diff := engine.DetectMaterialChange(req.Baseline.SourceClassification, current)
		violations = append(violations, configDriftViolations(diff)...)
	}

	result := deriveResult(violations)
	runtimeFP := HashObservation(obs)
	baselineFP := ""
	if req.Baseline != nil {
		baselineFP = req.Baseline.Fingerprint
	}

	receipt := buildReceipt(obs, result, violations, runtimeFP, baselineFP, obs.PolicyVersion, now, "")

	return schema.VerifyRuntimeResponse{
		Result:  result,
		Receipt: receipt,
	}, nil
}

func observationAsClassifyRequest(obs schema.RuntimeObservation, baseline schema.ClassifyRequest) schema.ClassifyRequest {
	current := baseline
	if obs.ModelID != "" {
		current.ModelID = obs.ModelID
	}
	if obs.SystemPromptHash != "" {
		current.SystemPromptHash = obs.SystemPromptHash
	}
	// Connected tools are configuration state; per-event invocations are policy-checked separately.
	if obs.EventType == schema.RuntimeEventDeploy || obs.EventType == schema.RuntimeEventConfigChange {
		if len(obs.Tools) > 0 {
			current.Tools = append([]string(nil), obs.Tools...)
		}
	}
	if len(obs.DataCategories) > 0 {
		current.DataTypes = append([]schema.DataType(nil), obs.DataCategories...)
	}
	if obs.Geography != "" {
		current.GeographicScope = obs.Geography
	}
	if obs.AutonomyLevel != "" {
		current.AutonomyLevel = obs.AutonomyLevel
	}
	return current
}

func evaluateEnvelope(obs schema.RuntimeObservation, policy schema.RuntimePolicy) []schema.RuntimeViolation {
	var violations []schema.RuntimeViolation

	if obs.ModelID != "" {
		if containsString(policy.ForbiddenModels, obs.ModelID) {
			violations = append(violations, violation(
				"forbidden_model", schema.ViolationCritical, "model_id",
				"not in forbidden list", obs.ModelID, schema.ImpactReclassify,
			))
		} else if len(policy.ApprovedModels) > 0 && !containsString(policy.ApprovedModels, obs.ModelID) {
			violations = append(violations, violation(
				"unapproved_model", schema.ViolationWarn, "model_id",
				strings.Join(policy.ApprovedModels, ", "), obs.ModelID, schema.ImpactReassure,
			))
		}
	}

	for _, tool := range obs.Tools {
		if containsString(policy.ForbiddenTools, tool) {
			violations = append(violations, violation(
				"forbidden_tool", schema.ViolationCritical, "tools",
				"not in forbidden list", tool, schema.ImpactReassure,
			))
		} else if len(policy.ApprovedTools) > 0 && !containsString(policy.ApprovedTools, tool) {
			violations = append(violations, violation(
				"unapproved_tool", schema.ViolationCritical, "tools",
				strings.Join(policy.ApprovedTools, ", "), tool, schema.ImpactReassure,
			))
		}
	}

	for _, dt := range obs.DataCategories {
		if containsDataType(policy.ForbiddenDataCategories, dt) {
			violations = append(violations, violation(
				"forbidden_data_category", schema.ViolationCritical, "data_categories",
				"not forbidden", string(dt), schema.ImpactReclassify,
			))
		} else if len(policy.ApprovedDataCategories) > 0 && !containsDataType(policy.ApprovedDataCategories, dt) {
			violations = append(violations, violation(
				"unapproved_data_category", schema.ViolationWarn, "data_categories",
				dataTypesString(policy.ApprovedDataCategories), string(dt), schema.ImpactReassure,
			))
		}
	}

	if obs.Geography != "" {
		if containsGeo(policy.ForbiddenGeography, obs.Geography) {
			violations = append(violations, violation(
				"forbidden_geography", schema.ViolationCritical, "geography",
				"not forbidden", string(obs.Geography), schema.ImpactReclassify,
			))
		} else if len(policy.ApprovedGeography) > 0 && !containsGeo(policy.ApprovedGeography, obs.Geography) {
			violations = append(violations, violation(
				"unapproved_geography", schema.ViolationWarn, "geography",
				geoString(policy.ApprovedGeography), string(obs.Geography), schema.ImpactReassure,
			))
		}
	}

	if obs.AutonomyLevel != "" && policy.MaxAutonomy != "" {
		if autonomyRank[obs.AutonomyLevel] > autonomyRank[policy.MaxAutonomy] {
			violations = append(violations, violation(
				"autonomy_exceeded", schema.ViolationCritical, "autonomy_level",
				string(policy.MaxAutonomy), string(obs.AutonomyLevel), schema.ImpactReclassify,
			))
		}
	}

	if policy.HumanApprovalRequired && !obs.HumanApprovalGranted &&
		(obs.EventType == schema.RuntimeEventToolCall || obs.EventType == schema.RuntimeEventInference) {
		violations = append(violations, violation(
			"missing_human_approval", schema.ViolationCritical, "human_approval_granted",
			"true", "false", schema.ImpactReassure,
		))
	}

	return violations
}

func configDriftViolations(diff schema.MaterialChangeResult) []schema.RuntimeViolation {
	if !diff.Material {
		return nil
	}

	return []schema.RuntimeViolation{{
		Code:             "config_drift",
		Severity:         schema.ViolationWarn,
		Field:            strings.Join(diff.ChangedFields, ", "),
		Expected:         diff.PriorFingerprint,
		Actual:           diff.CurrentFingerprint,
		ConformityImpact: diff.ConformityImpact,
	}}
}

func violation(code string, sev schema.RuntimeViolationSeverity, field, expected, actual string, impact schema.ConformityImpact) schema.RuntimeViolation {
	return schema.RuntimeViolation{
		Code:             code,
		Severity:         sev,
		Field:            field,
		Expected:         expected,
		Actual:           actual,
		ConformityImpact: impact,
	}
}

func deriveResult(violations []schema.RuntimeViolation) schema.VerificationResult {
	state := schema.ConformityGreen
	var failedControls []string
	recommended := "continue"
	riskDelta := "none"

	hasCritical := false
	hasReclassify := false
	hasReassure := false

	for _, v := range violations {
		if v.Severity == schema.ViolationCritical {
			hasCritical = true
		}
		switch v.ConformityImpact {
		case schema.ImpactReclassify:
			hasReclassify = true
		case schema.ImpactReassure:
			hasReassure = true
		}
		failedControls = append(failedControls, v.Code)
	}

	switch {
	case hasCritical:
		state = schema.ConformityRed
		recommended = "block_and_alert"
		riskDelta = "elevated"
	case hasReclassify || hasReassure:
		state = schema.ConformityAmber
		if recommended == "continue" {
			recommended = "reverify_or_reclassify"
		}
		if hasReclassify {
			riskDelta = "reclassify_required"
		} else {
			riskDelta = "reassure_required"
		}
	}

	summary := "Runtime observation is within the approved operating envelope."
	if len(violations) > 0 {
		summary = fmt.Sprintf("%d runtime violation(s) detected.", len(violations))
	}

	return schema.VerificationResult{
		ConformityState:   state,
		Violations:        violations,
		RiskDelta:         riskDelta,
		ControlsFailed:    failedControls,
		RecommendedAction: recommended,
		Summary:           summary,
	}
}

func buildReceipt(
	obs schema.RuntimeObservation,
	result schema.VerificationResult,
	violations []schema.RuntimeViolation,
	runtimeFP, baselineFP, policyVersion string,
	verifiedAt time.Time,
	previousHash string,
) schema.VerificationReceipt {
	receipt := schema.VerificationReceipt{
		VerificationID:      newVerificationID(obs, verifiedAt),
		SystemID:            obs.SystemID,
		VerifiedAt:          verifiedAt,
		ObservedAt:          obs.Timestamp,
		RuntimeFingerprint:  runtimeFP,
		BaselineFingerprint: baselineFP,
		PolicyVersion:       policyVersion,
		Violations:          violations,
		ConformityState:     result.ConformityState,
		PreviousReceiptHash: previousHash,
	}
	receipt.ReceiptHash = hashReceipt(receipt)
	return receipt
}

func newVerificationID(obs schema.RuntimeObservation, at time.Time) string {
	payload := fmt.Sprintf("%s:%s:%d", obs.SystemID, obs.DeploymentID, at.UnixNano())
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:16])
}

func hashReceipt(r schema.VerificationReceipt) string {
	stub := r
	stub.ReceiptHash = ""
	stub.Signature = ""
	stub.PublicKey = ""
	stub.Algorithm = ""
	data, err := json.Marshal(stub)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// HashObservation returns a stable fingerprint of runtime observation metadata.
func HashObservation(obs schema.RuntimeObservation) string {
	tools := append([]string(nil), obs.Tools...)
	sort.Strings(tools)

	dts := append([]schema.DataType(nil), obs.DataCategories...)
	sort.Slice(dts, func(i, j int) bool { return dts[i] < dts[j] })

	normalized := map[string]any{
		"system_id":               obs.SystemID,
		"deployment_id":           obs.DeploymentID,
		"model_id":                obs.ModelID,
		"system_prompt_hash":      obs.SystemPromptHash,
		"tools":                   tools,
		"event_type":              obs.EventType,
		"data_categories":         dts,
		"geography":               obs.Geography,
		"autonomy_level":          obs.AutonomyLevel,
		"human_approval_granted":  obs.HumanApprovalGranted,
		"policy_version":          obs.PolicyVersion,
		"ruleset_version":         obs.RulesetVersion,
		"timestamp":               obs.Timestamp.UTC().Format(time.RFC3339Nano),
	}

	data, err := json.Marshal(normalized)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func containsDataType(list []schema.DataType, v schema.DataType) bool {
	for _, d := range list {
		if d == v {
			return true
		}
	}
	return false
}

func containsGeo(list []schema.GeographicScope, v schema.GeographicScope) bool {
	for _, g := range list {
		if g == v {
			return true
		}
	}
	return false
}

func dataTypesString(dts []schema.DataType) string {
	parts := make([]string, len(dts))
	for i, d := range dts {
		parts[i] = string(d)
	}
	return strings.Join(parts, ", ")
}

func geoString(gs []schema.GeographicScope) string {
	parts := make([]string, len(gs))
	for i, g := range gs {
		parts[i] = string(g)
	}
	return strings.Join(parts, ", ")
}
