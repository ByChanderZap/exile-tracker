package frontend

import (
	"database/sql"
	"errors"
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
	// show characters by accound and search characters within an account
	router.Get("/accounts/{accountId}/characters", h.handleCharactersByAccount)
	router.Get("/accounts/{accountId}/characters/search", h.handleCharactersSearchByAccount)
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

func (h *Handler) handleCharactersByAccount(w http.ResponseWriter, r *http.Request) {
	accId := chi.URLParam(r, "accountId")

	cs, err := h.repository.GetCharactersByAccountId(accId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No characters to show", http.StatusNotFound)
			return
		}
		h.log.Error().Err(err).Msg("Query to get all characters by account id failed")
		http.Error(w, "Failed to load characters", http.StatusBadRequest)
		return
	}
	templates.CharactersByAccountId(cs, accId, utils.StringValue).Render(r.Context(), w)
}

func (h *Handler) handleCharactersSearchByAccount(w http.ResponseWriter, r *http.Request) {
	accId := chi.URLParam(r, "accountId")
	searchTerm := r.URL.Query().Get("q")

	var characters []models.Character
	var err error

	if searchTerm == "" {
		characters, err = h.repository.GetCharactersByAccountId(accId)
	} else {
		characters, err = h.repository.SearchCharactersInAccount(repository.SearchCharactersInAccountParams{
			AccountId: accId,
			Query:     searchTerm,
		})
	}

	if err != nil {
		h.log.Error().Err(err).Msg("Query to search characters in account failed")
		http.Error(w, "Failed to search characters by account", http.StatusInternalServerError)
		return
	}

	templates.CharactersTable(characters, utils.StringValue).Render(r.Context(), w)
}
