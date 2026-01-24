package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	var notFound *domain.ErrNotFound
	var validation *domain.ErrValidation
	var conflict *domain.ErrConflict

	switch {
	case errors.As(err, &notFound):
		respondJSON(w, http.StatusNotFound, errorBody{
			Error: errorDetail{Code: "not_found", Message: err.Error()},
		})
	case errors.As(err, &validation):
		respondJSON(w, http.StatusBadRequest, errorBody{
			Error: errorDetail{Code: "validation_error", Message: err.Error()},
		})
	case errors.As(err, &conflict):
		respondJSON(w, http.StatusConflict, errorBody{
			Error: errorDetail{Code: "conflict", Message: err.Error()},
		})
	default:
		slog.Error("unexpected error", "error", err)
		respondJSON(w, http.StatusInternalServerError, errorBody{
			Error: errorDetail{Code: "internal_error", Message: "internal server error"},
		})
	}
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	respondJSON(w, http.StatusBadRequest, errorBody{
		Error: errorDetail{Code: "bad_request", Message: msg},
	})
}
