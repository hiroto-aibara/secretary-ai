package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	"github.com/hiroto-aibara/secretary-ai/internal/handler"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
)

type mockBoardRepo struct {
	boards  []domain.Board
	board   *domain.Board
	getErr  error
	saveErr error
	delErr  error
}

func (m *mockBoardRepo) List(_ context.Context) ([]domain.Board, error) {
	return m.boards, nil
}

func (m *mockBoardRepo) Get(_ context.Context, id string) (*domain.Board, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.board != nil && m.board.ID == id {
		return m.board, nil
	}
	return nil, &domain.ErrNotFound{Resource: "board", ID: id}
}

func (m *mockBoardRepo) Save(_ context.Context, board *domain.Board) error {
	m.board = board
	return m.saveErr
}

func (m *mockBoardRepo) Delete(_ context.Context, _ string) error {
	return m.delErr
}

func newBoardRouter(repo *mockBoardRepo) *chi.Mux {
	uc := usecase.NewBoardUseCase(repo)
	h := handler.NewBoardHandler(uc)
	r := chi.NewRouter()
	h.Register(r)
	return r
}

func TestBoardHandler_List(t *testing.T) {
	repo := &mockBoardRepo{boards: []domain.Board{
		{ID: "a", Name: "A", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var boards []domain.Board
	if err := json.NewDecoder(w.Body).Decode(&boards); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(boards) != 1 {
		t.Errorf("got %d boards, want 1", len(boards))
	}
}

func TestBoardHandler_Create(t *testing.T) {
	repo := &mockBoardRepo{
		getErr: &domain.ErrNotFound{Resource: "board", ID: "new"},
	}
	r := newBoardRouter(repo)

	body := `{"id":"new","name":"New Board","lists":[{"id":"todo","name":"Todo"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestBoardHandler_Create_InvalidBody(t *testing.T) {
	repo := &mockBoardRepo{}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/boards", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBoardHandler_Get(t *testing.T) {
	repo := &mockBoardRepo{
		board: &domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/test", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestBoardHandler_Get_NotFound(t *testing.T) {
	repo := &mockBoardRepo{}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/missing", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestBoardHandler_Update(t *testing.T) {
	repo := &mockBoardRepo{
		board: &domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}
	r := newBoardRouter(repo)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/boards/test", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestBoardHandler_Delete(t *testing.T) {
	repo := &mockBoardRepo{
		board: &domain.Board{ID: "test", Name: "Test"},
	}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/test", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestBoardHandler_Update_InvalidBody(t *testing.T) {
	repo := &mockBoardRepo{
		board: &domain.Board{ID: "test", Name: "Test"},
	}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodPut, "/api/boards/test", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBoardHandler_Update_NotFound(t *testing.T) {
	repo := &mockBoardRepo{}
	r := newBoardRouter(repo)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/boards/missing", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestBoardHandler_Delete_NotFound(t *testing.T) {
	repo := &mockBoardRepo{}
	r := newBoardRouter(repo)

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/missing", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
