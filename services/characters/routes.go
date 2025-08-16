package characters

import (
	"net/http"

	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/snapshots/{accountName}/{characterName}", h.handleGetCharacterSnapshot)
}

func (h *Handler) handleGetCharacterSnapshot(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message":   "Character snapshot endpoint is not implemented yet",
		"account":   chi.URLParam(r, "accountName"),
		"character": chi.URLParam(r, "characterName"),
	})
}
