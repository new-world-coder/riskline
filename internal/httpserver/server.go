package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/runtime"
	"github.com/new-world-coder/riskline/pkg/schema"
)

// Server is a thin HTTP wrapper around the classification engine.
type Server struct {
	eng *engine.Engine
	mux *http.ServeMux
}

func New(eng *engine.Engine) *Server {
	s := &Server{eng: eng, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("POST /v1/classify", s.handleClassify)
	s.mux.HandleFunc("POST /v1/diff", s.handleDiff)
	s.mux.HandleFunc("POST /v1/assure", s.handleAssure)
	s.mux.HandleFunc("POST /v1/runtime/register", s.handleRuntimeRegister)
	s.mux.HandleFunc("POST /v1/runtime/verify", s.handleRuntimeVerify)
	s.mux.HandleFunc("POST /v1/runtime/observe", s.handleRuntimeObserve)
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleClassify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req schema.ClassifyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	resp, err := s.eng.Classify(req)
	if err != nil {
		if ve, ok := err.(*engine.ValidationError); ok {
			writeError(w, http.StatusBadRequest, "validation failed", ve.Details)
			return
		}
		writeError(w, http.StatusInternalServerError, "classification failed", []string{err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func (s *Server) handleDiff(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req schema.DiffRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	resp := engine.DetectMaterialChange(req.Baseline, req.Current)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func (s *Server) handleAssure(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req schema.AssureRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	resp := assure.Evaluate(req)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func (s *Server) handleRuntimeRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req schema.RegisterRuntimeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	resp, err := runtime.Register(req, timeNow())
	if err != nil {
		writeError(w, http.StatusBadRequest, "registration failed", []string{err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func (s *Server) handleRuntimeVerify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req schema.VerifyRuntimeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	resp, err := runtime.Verify(req, timeNow())
	if err != nil {
		writeError(w, http.StatusBadRequest, "verification failed", []string{err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
}

func (s *Server) handleRuntimeObserve(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var obs schema.RuntimeObservation
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&obs); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body", []string{err.Error()})
		return
	}

	// Stateless OSS API: accept and fingerprint only. Hosted persistence is riskline-cloud.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"accepted":    true,
		"message":     "observation accepted; use POST /v1/runtime/verify for conformity verdict",
		"fingerprint": runtime.HashObservation(obs),
	})
}

func timeNow() time.Time {
	return time.Now().UTC()
}

func writeError(w http.ResponseWriter, status int, msg string, details []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(schema.ErrorResponse{
		Error:      msg,
		Details:    details,
		Disclaimer: schema.Disclaimer,
	})
}
