package server

import (
	"context"
	"encoding/json"
	svcHealth "github.com/edwinfinlay/service-health-aggregator-nable/internal/app/service-health-aggregator"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(ctx context.Context, healthSvc *svcHealth.HealthService) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/aggregate", func(w http.ResponseWriter, r *http.Request) {
		resp, err := healthSvc.CheckHealth(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return &Server{
		mux: mux,
	}

}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8000", s.mux)
}
