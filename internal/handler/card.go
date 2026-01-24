package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
)

type CardHandler struct {
	uc *usecase.CardUseCase
}

func NewCardHandler(uc *usecase.CardUseCase) *CardHandler {
	return &CardHandler{uc: uc}
}

func (h *CardHandler) Register(r chi.Router) {
	r.Get("/api/boards/{id}/cards", h.list)
	r.Post("/api/boards/{id}/cards", h.create)
	r.Get("/api/boards/{id}/cards/{cardId}", h.get)
	r.Put("/api/boards/{id}/cards/{cardId}", h.update)
	r.Delete("/api/boards/{id}/cards/{cardId}", h.delete)
	r.Patch("/api/boards/{id}/cards/{cardId}/move", h.move)
	r.Patch("/api/boards/{id}/cards/{cardId}/archive", h.archive)
}

func (h *CardHandler) list(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	includeArchived := r.URL.Query().Get("archived") == "true"

	cards, err := h.uc.List(r.Context(), boardID, includeArchived)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, cards)
}

func (h *CardHandler) create(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")

	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	created, err := h.uc.Create(r.Context(), boardID, &card)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, created)
}

func (h *CardHandler) get(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	card, err := h.uc.Get(r.Context(), boardID, cardID)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, card)
}

func (h *CardHandler) update(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	var updates domain.Card
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	updated, err := h.uc.Update(r.Context(), boardID, cardID, &updates)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

func (h *CardHandler) delete(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	if err := h.uc.Delete(r.Context(), boardID, cardID); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type moveRequest struct {
	List  string `json:"list"`
	Order int    `json:"order"`
}

func (h *CardHandler) move(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	if req.List == "" {
		writeBadRequest(w, "list is required")
		return
	}

	card, err := h.uc.Move(r.Context(), boardID, cardID, req.List, req.Order)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, card)
}

type archiveRequest struct {
	Archived bool `json:"archived"`
}

func (h *CardHandler) archive(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "id")
	cardID := chi.URLParam(r, "cardId")

	var req archiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	card, err := h.uc.Archive(r.Context(), boardID, cardID)
	if err != nil {
		writeError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, card)
}
