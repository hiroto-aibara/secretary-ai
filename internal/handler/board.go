package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
)

type BoardHandler struct {
	uc *usecase.BoardUseCase
}

func NewBoardHandler(uc *usecase.BoardUseCase) *BoardHandler {
	return &BoardHandler{uc: uc}
}

func (h *BoardHandler) Register(r chi.Router) {
	r.Get("/api/boards", h.list)
	r.Post("/api/boards", h.create)
	r.Get("/api/boards/{id}", h.get)
	r.Put("/api/boards/{id}", h.update)
	r.Delete("/api/boards/{id}", h.delete)
}

func (h *BoardHandler) list(w http.ResponseWriter, r *http.Request) {
	boards, err := h.uc.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, boards)
}

func (h *BoardHandler) create(w http.ResponseWriter, r *http.Request) {
	var board domain.Board
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	created, err := h.uc.Create(r.Context(), &board)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, created)
}

func (h *BoardHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	board, err := h.uc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, board)
}

func (h *BoardHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var board domain.Board
	if err := json.NewDecoder(r.Body).Decode(&board); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	updated, err := h.uc.Update(r.Context(), id, &board)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

func (h *BoardHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
