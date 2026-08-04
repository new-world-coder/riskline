package schema

import (
	"encoding/json"
	"testing"
)

func FuzzDecodeClassifyRequest(f *testing.F) {
	seed, _ := json.Marshal(ClassifyRequest{
		Purpose:            "Screen job applicants",
		DataTypes:          []DataType{DataEmployment},
		DeploymentContext:  DeploySaaSB2B,
		AutonomyLevel:      AutonomyDecisionSupport,
		AffectedPopulation: PopJobApplicants,
	})
	f.Add(seed)
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"purpose":"x","data_types":["personal_data"],"deployment_context":"saas_b2b","autonomy_level":"content_generation","affected_population":"customers"}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var req ClassifyRequest
		_ = json.Unmarshal(data, &req)
	})
}
