package api

import (
	"context"
	"net/http"

	"github.com/ByChanderZap/exile-tracker/services/characters"
	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	addr   string
	server *http.Server
}

func NewAPIServer(addr string) *APIServer {
	utils.BaseLogger.Info().Msg(addr)
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Start() error {
	log := utils.ChildLogger("api")

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	// character endpoints
	cHandler := characters.NewHandler()
	cHandler.RegisterRoutes(v1Router)

	router.Mount("/api/v1", v1Router)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Info().Msgf("Starting API server, listening on %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *APIServer) Stop(ctx context.Context) error {
	log := utils.ChildLogger("api")
	log.Info().Msg("Stopping API server")
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
