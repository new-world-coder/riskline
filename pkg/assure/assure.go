package assure

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/new-world-coder/riskline/pkg/schema"
)

const hashAlgorithm = "sha256"

// Evaluate verifies technical controls from a classification against probe results.
// Probes are keyed by technical_hook first, then control id as fallback.
func Evaluate(req schema.AssureRequest) schema.AssureResponse {
	controls := req.Classification.TechnicalControls
	if len(controls) == 0 {
		return schema.AssureResponse{
			ConformityState: schema.ConformityGreen,
			Summary:         "No technical controls required for this classification; conformity state defaults to green (obligations only).",
			Disclaimer:      schema.Disclaimer,
		}
	}

	var verdicts []schema.ControlVerdict
	var records []schema.EvidenceRecord
	failed := 0
	unverified := 0
	prevHash := req.PreviousHash

	for _, tc := range controls {
		passed, reason, known := probeResult(req.Probes, tc)
		v := schema.ControlVerdict{
			ControlID:     tc.ID,
			TechnicalHook: tc.TechnicalHook,
			PaperRef:      tc.PaperRef,
			EvidenceType:  tc.EvidenceType,
			Passed:        passed,
			Reason:        reason,
		}
		verdicts = append(verdicts, v)

		if !known {
			unverified++
			continue
		}

		result := "pass"
		if !passed {
			result = "fail"
			failed++
		}

		rec := buildEvidenceRecord(req.Classification, tc, result, prevHash)
		records = append(records, rec)
		prevHash = rec.ContentHash
	}

	state, summary := deriveState(failed, unverified, len(controls))

	return schema.AssureResponse{
		ConformityState: state,
		ControlVerdicts: verdicts,
		EvidenceRecords: records,
		UnverifiedCount: unverified,
		FailedCount:     failed,
		Summary:         summary,
		Disclaimer:      schema.Disclaimer,
	}
}

func probeResult(probes map[string]bool, tc schema.TechnicalControl) (passed bool, reason string, known bool) {
	if probes == nil {
		return false, "no probe results supplied", false
	}
	if v, ok := probes[tc.TechnicalHook]; ok {
		if v {
			return true, "probe passed", true
		}
		return false, "probe reported failure", true
	}
	if v, ok := probes[tc.ID]; ok {
		if v {
			return true, "probe passed (by control id)", true
		}
		return false, "probe reported failure (by control id)", true
	}
	return false, "no probe registered for hook or control id", false
}

func deriveState(failed, unverified, total int) (schema.ConformityState, string) {
	switch {
	case failed > 0:
		return schema.ConformityRed,
			fmt.Sprintf("%d of %d controls failed verification; prior conformity evidence may be invalid.", failed, total)
	case unverified > 0:
		return schema.ConformityAmber,
			fmt.Sprintf("%d of %d controls have no probe result; verification incomplete.", unverified, total)
	default:
		return schema.ConformityGreen,
			fmt.Sprintf("All %d required controls verified.", total)
	}
}

func buildEvidenceRecord(classification schema.ClassifyResponse, tc schema.TechnicalControl, result, previousHash string) schema.EvidenceRecord {
	payload := map[string]string{
		"system_name":     classification.Name,
		"ruleset_version": classification.RulesetVersion,
		"control_id":      tc.ID,
		"technical_hook":  tc.TechnicalHook,
		"paper_ref":       tc.PaperRef,
		"evidence_type":   tc.EvidenceType,
		"result":          result,
		"previous_hash":   previousHash,
	}
	hash := canonicalHash(payload)
	return schema.EvidenceRecord{
		ID:             tc.ID + "@" + classification.RulesetVersion,
		Timestamp:      time.Now().UTC(),
		SystemName:     classification.Name,
		RulesetVersion: classification.RulesetVersion,
		ControlID:      tc.ID,
		Result:         result,
		ContentHash:    hash,
		PreviousHash:   previousHash,
		Algorithm:      hashAlgorithm,
	}
}

// HashClassification returns a stable fingerprint of the classification input
// fields that affect tier assignment (for material-change detection).
func HashClassification(req schema.ClassifyRequest) string {
	normalized := normalizeRequest(req)
	return canonicalHash(normalized)
}

func canonicalHash(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	// Stable key order for maps is handled by using struct or sorted keys in normalize
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func normalizeRequest(req schema.ClassifyRequest) map[string]any {
	dts := append([]schema.DataType(nil), req.DataTypes...)
	sort.Slice(dts, func(i, j int) bool { return dts[i] < dts[j] })

	regs := append([]string(nil), req.Regimes...)
	sort.Strings(regs)

	tools := append([]string(nil), req.Tools...)
	sort.Strings(tools)

	return map[string]any{
		"name":                                      req.Name,
		"purpose":                                   strings.TrimSpace(strings.ToLower(req.Purpose)),
		"data_types":                                dts,
		"deployment_context":                        req.DeploymentContext,
		"autonomy_level":                            req.AutonomyLevel,
		"affected_population":                       req.AffectedPopulation,
		"geographic_scope":                          req.GeographicScope,
		"regimes":                                   regs,
		"uses_biometric_categorisation":             req.UsesBiometricCategorisation,
		"real_time_remote_biometric_id":             req.RealTimeRemoteBiometricID,
		"social_scoring":                            req.SocialScoring,
		"emotion_recognition_workplace_or_education": req.EmotionRecognitionWorkplaceOrEducation,
		"manipulative_techniques":                   req.ManipulativeTechniques,
		"exploits_vulnerabilities":                  req.ExploitsVulnerabilities,

		// Runtime metadata
		"model_id":              req.ModelID,
		"system_prompt_hash":    req.SystemPromptHash,
		"tools":                 tools,
		"human_approval_required": req.HumanApprovalRequired,
	}
}

// VerifyEvidenceChain checks that each record's previous_hash matches the prior
// record's content_hash (local integrity chain, not blockchain).
func VerifyEvidenceChain(records []schema.EvidenceRecord) (bool, string) {
	var prev string
	for i, rec := range records {
		if i > 0 && rec.PreviousHash != prev {
			return false, fmt.Sprintf("chain break at record %d: expected previous_hash %s, got %s", i, prev, rec.PreviousHash)
		}
		prev = rec.ContentHash
	}
	return true, "evidence chain intact"
}
