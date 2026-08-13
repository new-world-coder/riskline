package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/schema"
)

func TestClassifyEndpoint(t *testing.T) {
	eng, err := engine.Default()
	if err != nil {
		t.Fatal(err)
	}
	srv := New(eng)

	body := `{
		"name": "Hiring Assist",
		"purpose": "Screen job applicants and rank candidates for interview",
		"data_types": ["personal_data", "employment_data"],
		"deployment_context": "saas_b2b",
		"autonomy_level": "decision_support",
		"affected_population": "job_applicants"
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/classify", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var resp schema.ClassifyResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.RiskTier != schema.TierHighRisk {
		t.Fatalf("tier=%s", resp.RiskTier)
	}
	if resp.Disclaimer == "" {
		t.Fatal("disclaimer missing from API response")
	}
}

func TestClassifyBadJSON(t *testing.T) {
	eng, err := engine.Default()
	if err != nil {
		t.Fatal(err)
	}
	srv := New(eng)
	req := httptest.NewRequest(http.MethodPost, "/v1/classify", bytes.NewBufferString(`{`))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rr.Code)
	}
	var er schema.ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &er); err != nil {
		t.Fatal(err)
	}
	if er.Disclaimer == "" {
		t.Fatal("disclaimer missing from error response")
	}
}

func TestClassifyUnknownRegime(t *testing.T) {
	eng, err := engine.Default()
	if err != nil {
		t.Fatal(err)
	}
	srv := New(eng)
	body := `{
		"purpose": "toy",
		"data_types": ["other"],
		"deployment_context": "other",
		"autonomy_level": "content_generation",
		"affected_population": "other",
		"regimes": ["mas-feat"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/classify", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHealthz(t *testing.T) {
	eng, err := engine.Default()
	if err != nil {
		t.Fatal(err)
	}
	srv := New(eng)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 || rr.Body.String() != "ok" {
		t.Fatalf("healthz failed: %d %s", rr.Code, rr.Body.String())
	}
}
