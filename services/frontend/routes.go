package frontend

import (
	"net/http"

	"github.com/ByChanderZap/exile-tracker/cmd/web/templates"
	"github.com/ByChanderZap/exile-tracker/models"
	"github.com/ByChanderZap/exile-tracker/repository"
	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	repository *repository.Repository
	log        zerolog.Logger
}

func NewHandler(db *repository.Repository, logger zerolog.Logger) *Handler {
	return &Handler{
		repository: db,
		log:        logger,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Use(utils.ZerologMiddleware(h.log))
	fs := http.FileServer(http.Dir("cmd/web/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", h.handleHomePage)
	router.Get("/search", h.handleSearchAccounts)
	// router.Get("/accounts", h.handleAccounts)
}

func (h *Handler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	// templ.Handler(templates.Main("alex"))

	accounts, err := h.repository.GetAllAccounts()
	if err != nil {
		h.log.Error().Err(err).Msg("Query to get accounts failed")
		http.Error(w, "Filed to load accounts", http.StatusInternalServerError)
		return
	}

	templates.Main(accounts, utils.StringValue).Render(r.Context(), w)
}

func (h *Handler) handleSearchAccounts(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")

	var accounts []models.Account
	var err error

	if searchTerm == "" {
		accounts, err = h.repository.GetAllAccounts()
	} else {
		accounts, err = h.repository.SearchAccounts(searchTerm)
	}

	if err != nil {
		h.log.Error().Err(err).Msg("Query to search accounts failed")
		http.Error(w, "Failed to search accounts", http.StatusInternalServerError)
		return
	}

	templates.AccountsTable(accounts, utils.StringValue).Render(r.Context(), w)
}
