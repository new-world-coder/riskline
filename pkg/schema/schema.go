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
	Name                                    string             `json:"name,omitempty" yaml:"name,omitempty"`
	Purpose                                 string             `json:"purpose" yaml:"purpose"`
	DataTypes                               []DataType         `json:"data_types" yaml:"data_types"`
	DeploymentContext                       DeploymentContext  `json:"deployment_context" yaml:"deployment_context"`
	AutonomyLevel                           AutonomyLevel      `json:"autonomy_level" yaml:"autonomy_level"`
	AffectedPopulation                      AffectedPopulation `json:"affected_population" yaml:"affected_population"`
	GeographicScope                         GeographicScope    `json:"geographic_scope,omitempty" yaml:"geographic_scope,omitempty"`
	UsesBiometricCategorisation             bool               `json:"uses_biometric_categorisation,omitempty" yaml:"uses_biometric_categorisation,omitempty"`
	RealTimeRemoteBiometricID               bool               `json:"real_time_remote_biometric_id,omitempty" yaml:"real_time_remote_biometric_id,omitempty"`
	SocialScoring                           bool               `json:"social_scoring,omitempty" yaml:"social_scoring,omitempty"`
	EmotionRecognitionWorkplaceOrEducation  bool               `json:"emotion_recognition_workplace_or_education,omitempty" yaml:"emotion_recognition_workplace_or_education,omitempty"`
	ManipulativeTechniques                  bool               `json:"manipulative_techniques,omitempty" yaml:"manipulative_techniques,omitempty"`
	ExploitsVulnerabilities                 bool               `json:"exploits_vulnerabilities,omitempty" yaml:"exploits_vulnerabilities,omitempty"`
}

type MatchedRule struct {
	ID              string   `json:"id"`
	Tier            RiskTier `json:"tier"`
	ArticleOrAnnex  string   `json:"article_or_annex"`
	Summary         string   `json:"summary"`
	RulesetVersion  string   `json:"ruleset_version"`
	LastUpdated     string   `json:"last_updated"`
}

type ClassifyResponse struct {
	Name                 string        `json:"name,omitempty"`
	RiskTier             RiskTier      `json:"risk_tier"`
	RulesetVersion       string        `json:"ruleset_version"`
	LastUpdated          string        `json:"last_updated"`
	MatchedRules         []MatchedRule `json:"matched_rules"`
	Rationale            string        `json:"rationale"`
	RecommendedControls  []string      `json:"recommended_controls"`
	JudgmentCalls        []string      `json:"judgment_calls,omitempty"`
	Disclaimer           string        `json:"disclaimer"`
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
