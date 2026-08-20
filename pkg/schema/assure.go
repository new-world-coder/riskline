package schema

import "time"

// ConformityState is the continuously computed compliance posture for a
// deployed system relative to its classified obligations.
type ConformityState string

const (
	ConformityGreen ConformityState = "green" // all required controls verified
	ConformityAmber ConformityState = "amber" // material change or partial verification
	ConformityRed   ConformityState = "red"   // control failure or missing evidence
)

// ConformityImpact describes what a detected change requires.
type ConformityImpact string

const (
	ImpactNone       ConformityImpact = "none"
	ImpactReassure   ConformityImpact = "reassure"   // tier unchanged; re-run probes
	ImpactReclassify ConformityImpact = "reclassify" // system description changed materially
)

// MaterialChangeResult compares two system descriptions (or fingerprints).
type MaterialChangeResult struct {
	Material          bool             `json:"material"`
	ConformityImpact  ConformityImpact `json:"conformity_impact"`
	ChangedFields     []string         `json:"changed_fields,omitempty"`
	PriorFingerprint  string           `json:"prior_fingerprint"`
	CurrentFingerprint string          `json:"current_fingerprint"`
	Summary           string           `json:"summary"`
	Disclaimer        string           `json:"disclaimer"`
}

// DiffRequest compares baseline and current system descriptions.
type DiffRequest struct {
	Baseline ClassifyRequest `json:"baseline"`
	Current  ClassifyRequest `json:"current"`
}

// ControlVerdict is the outcome of verifying one technical control hook.
type ControlVerdict struct {
	ControlID     string `json:"control_id"`
	TechnicalHook string `json:"technical_hook"`
	PaperRef      string `json:"paper_ref"`
	EvidenceType  string `json:"evidence_type"`
	Passed        bool   `json:"passed"`
	Reason        string `json:"reason,omitempty"`
}

// EvidenceRecord is a tamper-evident audit artifact. Hashes are computed locally
// (SHA-256 over canonical JSON). This is not a public blockchain ledger.
type EvidenceRecord struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	SystemName     string    `json:"system_name,omitempty"`
	RulesetVersion string    `json:"ruleset_version"`
	ControlID      string    `json:"control_id"`
	Result         string    `json:"result"` // pass | fail | unknown
	ContentHash    string    `json:"content_hash"`
	PreviousHash   string    `json:"previous_hash,omitempty"`
	Algorithm      string    `json:"algorithm"` // e.g. sha256
}

// AssureRequest verifies classified obligations against probe results.
// Probes map technical_hook (or control id) to pass/fail from customer systems.
type AssureRequest struct {
	Classification ClassifyResponse `json:"classification"`
	Probes         map[string]bool  `json:"probes"`
	PreviousHash   string           `json:"previous_hash,omitempty"`
}

// AssureResponse is the Layer 2 Assure output.
type AssureResponse struct {
	ConformityState  ConformityState    `json:"conformity_state"`
	ControlVerdicts  []ControlVerdict   `json:"control_verdicts"`
	EvidenceRecords  []EvidenceRecord   `json:"evidence_records"`
	UnverifiedCount  int                `json:"unverified_count"`
	FailedCount      int                `json:"failed_count"`
	Summary          string             `json:"summary"`
	Disclaimer       string             `json:"disclaimer"`
}

// EvidenceBundle is a locally signed assurance export (Ed25519 over canonical
// JSON of payload). Not a public blockchain ledger — verify offline with the
// embedded public_key.
type EvidenceBundle struct {
	Payload    AssureResponse `json:"payload"`
	Signature  string         `json:"signature"`
	PublicKey  string         `json:"public_key"`
	Algorithm  string         `json:"algorithm"`
	Disclaimer string         `json:"disclaimer"`
}
