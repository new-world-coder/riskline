package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/new-world-coder/riskline/pkg/engine"
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

func writeError(w http.ResponseWriter, status int, msg string, details []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(schema.ErrorResponse{
		Error:      msg,
		Details:    details,
		Disclaimer: schema.Disclaimer,
	})
}
