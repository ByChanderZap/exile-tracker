package pobsnapshots

import (
	"net/http"

	"github.com/ByChanderZap/exile-tracker/repository"
	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repository *repository.Repository
}

func NewHandler(db *repository.Repository) *Handler {
	return &Handler{
		repository: db,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/pobsnapshots/character/{characterId}", h.handleGetSnapshotsByCharacter)
	router.Get("/pobsnapshots/character/{characterId}/latest", h.handleGetLatestSnapshot)
	router.Get("/pobsnapshots/{id}", h.handleGetSnapshotByID)
}

func (h *Handler) handleGetSnapshotsByCharacter(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	snapshots, err := h.repository.GetSnapshotsByCharacter(characterId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, snapshots)
}

func (h *Handler) handleGetLatestSnapshot(w http.ResponseWriter, r *http.Request) {
	characterId := chi.URLParam(r, "characterId")
	snapshot, err := h.repository.GetLatestSnapshotByCharacter(characterId)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, snapshot)
}

func (h *Handler) handleGetSnapshotByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	snapshot, err := h.repository.GetSnapshotByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, snapshot)
}
