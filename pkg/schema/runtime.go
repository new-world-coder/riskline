package schema

import "time"

// RuntimeEventType categorizes runtime signals submitted by deployed systems.
type RuntimeEventType string

const (
	RuntimeEventInference     RuntimeEventType = "inference"
	RuntimeEventToolCall      RuntimeEventType = "tool_call"
	RuntimeEventHumanOverride RuntimeEventType = "human_override"
	RuntimeEventDeploy        RuntimeEventType = "deploy"
	RuntimeEventConfigChange  RuntimeEventType = "config_change"
)

// RuntimeObservation is a metadata-only runtime signal from a deployed system.
type RuntimeObservation struct {
	SystemID              string            `json:"system_id"`
	DeploymentID          string            `json:"deployment_id"`
	ModelID               string            `json:"model_id,omitempty"`
	SystemPromptHash      string            `json:"system_prompt_hash,omitempty"`
	Tools                 []string          `json:"tools,omitempty"`
	EventType             RuntimeEventType  `json:"event_type"`
	DataCategories        []DataType        `json:"data_categories,omitempty"`
	Geography             GeographicScope   `json:"geography,omitempty"`
	AutonomyLevel         AutonomyLevel     `json:"autonomy_level,omitempty"`
	HumanApprovalGranted  bool              `json:"human_approval_granted"`
	PolicyVersion         string            `json:"policy_version,omitempty"`
	RulesetVersion        string            `json:"ruleset_version,omitempty"`
	Timestamp             time.Time         `json:"timestamp"`
	Metrics               map[string]float64 `json:"metrics,omitempty"`
}

// RuntimeBaseline captures the approved runtime posture at registration time.
type RuntimeBaseline struct {
	SystemID             string          `json:"system_id"`
	DeploymentID         string          `json:"deployment_id"`
	Fingerprint          string          `json:"fingerprint"`
	ApprovedAt           time.Time       `json:"approved_at"`
	SourceClassification ClassifyRequest `json:"source_classification"`
}

// RuntimePolicy defines approved and forbidden runtime boundaries for a deployment.
type RuntimePolicy struct {
	ApprovedModels            []string          `json:"approved_models,omitempty"`
	ForbiddenModels           []string          `json:"forbidden_models,omitempty"`
	ApprovedTools             []string          `json:"approved_tools,omitempty"`
	ForbiddenTools            []string          `json:"forbidden_tools,omitempty"`
	ApprovedDataCategories    []DataType        `json:"approved_data_categories,omitempty"`
	ForbiddenDataCategories   []DataType        `json:"forbidden_data_categories,omitempty"`
	ApprovedGeography         []GeographicScope `json:"approved_geography,omitempty"`
	ForbiddenGeography        []GeographicScope `json:"forbidden_geography,omitempty"`
	MaxAutonomy               AutonomyLevel     `json:"max_autonomy"`
	RequiredControls          []string          `json:"required_controls"`
	RiskTier                  RiskTier          `json:"risk_tier"`
	HumanApprovalRequired     bool              `json:"human_approval_required"`
}

// RuntimeViolationSeverity indicates how urgently a policy deviation should be handled.
type RuntimeViolationSeverity string

const (
	ViolationInfo     RuntimeViolationSeverity = "info"
	ViolationWarn     RuntimeViolationSeverity = "warn"
	ViolationCritical RuntimeViolationSeverity = "critical"
)

// RuntimeViolation describes one field-level deviation from the registered policy.
type RuntimeViolation struct {
	Code              string                   `json:"code"`
	Severity          RuntimeViolationSeverity `json:"severity"`
	Field             string                   `json:"field"`
	Expected          string                   `json:"expected"`
	Actual            string                   `json:"actual"`
	ConformityImpact  ConformityImpact         `json:"conformity_impact"`
}

// VerificationResult is the outcome of comparing a runtime observation to policy.
type VerificationResult struct {
	ConformityState    ConformityState    `json:"conformity_state"`
	Violations         []RuntimeViolation `json:"violations,omitempty"`
	RiskDelta          string             `json:"risk_delta,omitempty"`
	ControlsFailed     []string           `json:"controls_failed,omitempty"`
	RecommendedAction  string             `json:"recommended_action,omitempty"`
	Summary            string             `json:"summary"`
}

// VerificationReceipt is a tamper-evident runtime verification artifact.
// Hashes are computed locally (SHA-256 over canonical JSON). Signature fields
// follow the same local Ed25519 pattern as EvidenceBundle — not a public ledger.
type VerificationReceipt struct {
	VerificationID       string             `json:"verification_id"`
	SystemID               string             `json:"system_id"`
	VerifiedAt             time.Time          `json:"verified_at"`
	ObservedAt             time.Time          `json:"observed_at,omitempty"`
	RuntimeFingerprint     string             `json:"runtime_fingerprint"`
	BaselineFingerprint  string             `json:"baseline_fingerprint"`
	PolicyVersion          string             `json:"policy_version"`
	Violations             []RuntimeViolation `json:"violations,omitempty"`
	ConformityState        ConformityState    `json:"conformity_state"`
	ReceiptHash            string             `json:"receipt_hash"`
	PreviousReceiptHash    string             `json:"previous_receipt_hash,omitempty"`
	Signature              string             `json:"signature,omitempty"`
	PublicKey              string             `json:"public_key,omitempty"`
	Algorithm              string             `json:"algorithm,omitempty"`
}

// RegisterRuntimeRequest binds classification and assurance to a deployment baseline.
type RegisterRuntimeRequest struct {
	SystemID            string           `json:"system_id"`
	DeploymentID        string           `json:"deployment_id"`
	SourceDescription   ClassifyRequest  `json:"source_description"`
	Classification      ClassifyResponse `json:"classification"`
	Assure              AssureResponse   `json:"assure"`
	Probes              map[string]bool  `json:"probes,omitempty"`
}

// RegisterRuntimeResponse returns the derived baseline and policy envelope.
type RegisterRuntimeResponse struct {
	Baseline   RuntimeBaseline `json:"baseline"`
	Policy     RuntimePolicy   `json:"policy"`
	Disclaimer string          `json:"disclaimer"`
}

// VerifyRuntimeRequest compares an observation against registered or inline policy.
type VerifyRuntimeRequest struct {
	Observation         RuntimeObservation `json:"observation"`
	Policy              *RuntimePolicy     `json:"policy,omitempty"`
	Baseline            *RuntimeBaseline   `json:"baseline,omitempty"`
	PreviousReceiptHash string             `json:"previous_receipt_hash,omitempty"`
}

// VerificationReceiptBundle is a locally signed runtime verification export.
type VerificationReceiptBundle struct {
	Payload    VerificationReceipt `json:"payload"`
	Signature  string              `json:"signature"`
	PublicKey  string              `json:"public_key"`
	Algorithm  string              `json:"algorithm"`
	Disclaimer string              `json:"disclaimer"`
}

// VerifyRuntimeResponse returns the verification outcome and receipt stub.
type VerifyRuntimeResponse struct {
	Result  VerificationResult  `json:"result"`
	Receipt VerificationReceipt `json:"receipt"`
}
