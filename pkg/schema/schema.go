package schema

import "time"

// Disclaimer is attached to every successful and error API response.
const Disclaimer = "This classification is an advisory tool based on a versioned ruleset. It is not legal advice and is not a substitute for qualified counsel or a formal conformity assessment."

type DataType string

const (
	DataPersonal            DataType = "personal_data"
	DataSpecialCategory     DataType = "special_category_data"
	DataBiometric           DataType = "biometric_data"
	DataEmployment          DataType = "employment_data"
	DataFinancialCredit     DataType = "financial_credit_data"
	DataEducation           DataType = "education_data"
	DataHealth              DataType = "health_data"
	DataLawEnforcement      DataType = "law_enforcement_data"
	DataPubliclyAvailable   DataType = "publicly_available"
	DataSyntheticNonPersonal DataType = "synthetic_or_non_personal"
	DataOther               DataType = "other"
)

type DeploymentContext string

const (
	DeploySaaSB2B         DeploymentContext = "saas_b2b"
	DeploySaaSB2C         DeploymentContext = "saas_b2c"
	DeployOnPrem          DeploymentContext = "on_prem"
	DeployEmbedded        DeploymentContext = "embedded_product"
	DeployPublicAuthority DeploymentContext = "public_authority"
	DeployLawEnforcement  DeploymentContext = "law_enforcement"
	DeployInternalOnly    DeploymentContext = "internal_only"
	DeployOther           DeploymentContext = "other"
)

type AutonomyLevel string

const (
	AutonomyContentGeneration AutonomyLevel = "content_generation"
	AutonomyDecisionSupport   AutonomyLevel = "decision_support"
	AutonomyAutomatedDecision AutonomyLevel = "automated_decision"
	AutonomyAutonomousAction  AutonomyLevel = "autonomous_action"
)

type AffectedPopulation string

const (
	PopGeneralPublic     AffectedPopulation = "general_public"
	PopCustomers         AffectedPopulation = "customers"
	PopEmployees         AffectedPopulation = "employees"
	PopJobApplicants     AffectedPopulation = "job_applicants"
	PopStudents          AffectedPopulation = "students"
	PopPatients          AffectedPopulation = "patients"
	PopCreditConsumers   AffectedPopulation = "credit_consumers"
	PopSuspectsOrAccused AffectedPopulation = "suspects_or_accused"
	PopChildren          AffectedPopulation = "children"
	PopOther             AffectedPopulation = "other"
)

type GeographicScope string

const (
	GeoEU          GeographicScope = "eu"
	GeoEUAndGlobal GeographicScope = "eu_and_global"
	GeoNonEU       GeographicScope = "non_eu"
	GeoUnknown     GeographicScope = "unknown"
)

type RiskTier string

const (
	TierProhibited   RiskTier = "prohibited"
	TierHighRisk     RiskTier = "high_risk"
	TierLimitedRisk  RiskTier = "limited_risk"
	TierMinimalRisk  RiskTier = "minimal_risk"
)

// ClassifyRequest is the public input shape for classification.
type ClassifyRequest struct {
	Name                                   string             `json:"name,omitempty" yaml:"name,omitempty"`
	Purpose                                string             `json:"purpose" yaml:"purpose"`
	DataTypes                              []DataType         `json:"data_types" yaml:"data_types"`
	DeploymentContext                      DeploymentContext  `json:"deployment_context" yaml:"deployment_context"`
	AutonomyLevel                          AutonomyLevel      `json:"autonomy_level" yaml:"autonomy_level"`
	AffectedPopulation                     AffectedPopulation `json:"affected_population" yaml:"affected_population"`
	GeographicScope                        GeographicScope    `json:"geographic_scope,omitempty" yaml:"geographic_scope,omitempty"`
	// Regimes selects which legal/policy packs to evaluate. Empty means the
	// engine/CLI/API default (typically ["eu-ai-act"]). Not the same as
	// geographic_scope — do not infer regimes from geography.
	Regimes                                []string           `json:"regimes,omitempty" yaml:"regimes,omitempty"`
	UsesBiometricCategorisation            bool               `json:"uses_biometric_categorisation,omitempty" yaml:"uses_biometric_categorisation,omitempty"`
	RealTimeRemoteBiometricID              bool               `json:"real_time_remote_biometric_id,omitempty" yaml:"real_time_remote_biometric_id,omitempty"`
	SocialScoring                          bool               `json:"social_scoring,omitempty" yaml:"social_scoring,omitempty"`
	EmotionRecognitionWorkplaceOrEducation bool               `json:"emotion_recognition_workplace_or_education,omitempty" yaml:"emotion_recognition_workplace_or_education,omitempty"`
	ManipulativeTechniques                 bool               `json:"manipulative_techniques,omitempty" yaml:"manipulative_techniques,omitempty"`
	ExploitsVulnerabilities                bool               `json:"exploits_vulnerabilities,omitempty" yaml:"exploits_vulnerabilities,omitempty"`
}

type MatchedRule struct {
	ID              string   `json:"id"`
	Tier            RiskTier `json:"tier"`
	ArticleOrAnnex  string   `json:"article_or_annex"`
	Summary         string   `json:"summary"`
	RulesetVersion  string   `json:"ruleset_version"`
	LastUpdated     string   `json:"last_updated"`
}

// RegimeClassification is one pack's outcome. Kept separate from EU-shaped
// top-level fields so other jurisdictions are not forced into risk_tier enums.
type RegimeClassification struct {
	Regime              string        `json:"regime"`
	Character           string        `json:"character"`
	RiskTier            RiskTier      `json:"risk_tier"`
	RulesetVersion      string        `json:"ruleset_version"`
	LastUpdated         string        `json:"last_updated"`
	MatchedRules        []MatchedRule `json:"matched_rules"`
	Rationale           string        `json:"rationale"`
	RecommendedControls []string      `json:"recommended_controls"`
	JudgmentCalls       []string      `json:"judgment_calls,omitempty"`
}

type ClassifyResponse struct {
	Name                string        `json:"name,omitempty"`
	RiskTier            RiskTier      `json:"risk_tier"`
	RulesetVersion      string        `json:"ruleset_version"`
	LastUpdated         string        `json:"last_updated"`
	MatchedRules        []MatchedRule `json:"matched_rules"`
	Rationale           string        `json:"rationale"`
	RecommendedControls []string      `json:"recommended_controls"`
	JudgmentCalls       []string      `json:"judgment_calls,omitempty"`
	Disclaimer          string        `json:"disclaimer"`
	// Regime identifies which pack produced the top-level EU-compatible fields.
	// Omitted on the historical single-pack default path so golden JSON stays stable.
	Regime string `json:"regime,omitempty"`
	// Classifications is present when more than one regime was evaluated.
	// Single default eu-ai-act responses omit this for backward compatibility.
	Classifications []RegimeClassification `json:"classifications,omitempty"`
}

type ErrorResponse struct {
	Error      string   `json:"error"`
	Details    []string `json:"details,omitempty"`
	Disclaimer string   `json:"disclaimer"`
}

// ParseDate is a tiny helper so callers don't re-implement date formatting.
func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
