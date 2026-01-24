package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	"github.com/hiroto-aibara/secretary-ai/internal/handler"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
)

func TestWriteError_NotFound(t *testing.T) {
	repo := &mockBoardRepo{}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/nonexistent", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Error.Code != "not_found" {
		t.Errorf("code = %s, want not_found", body.Error.Code)
	}
}

func TestWriteError_Conflict(t *testing.T) {
	repo := &mockBoardRepo{
		board: &domain.Board{ID: "existing", Name: "Existing", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}

	uc := usecase.NewBoardUseCase(repo)
	h := handler.NewBoardHandler(uc)
	r := chi.NewRouter()
	h.Register(r)

	body := `{"id":"existing","name":"Conflict","lists":[{"id":"todo","name":"Todo"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusConflict, w.Body.String())
	}
}

func TestWriteError_InternalError(t *testing.T) {
	repo := &mockBoardRepo{
		board:   &domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
		saveErr: fmt.Errorf("disk full"),
	}

	uc := usecase.NewBoardUseCase(repo)
	h := handler.NewBoardHandler(uc)
	r := chi.NewRouter()
	h.Register(r)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/boards/test", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var errBody struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(w.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody.Error.Code != "internal_error" {
		t.Errorf("code = %s, want internal_error", errBody.Error.Code)
	}
}
