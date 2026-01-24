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

type mockCardRepo struct {
	cards     []domain.Card
	card      *domain.Card
	getErr    error
	saveErr   error
	delErr    error
	nextID    string
	nextIDErr error
	createErr error
}

func (m *mockCardRepo) ListByBoard(_ context.Context, _ string, _ bool) ([]domain.Card, error) {
	return m.cards, nil
}

func (m *mockCardRepo) Get(_ context.Context, _, cardID string) (*domain.Card, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.card != nil && m.card.ID == cardID {
		return m.card, nil
	}
	return nil, &domain.ErrNotFound{Resource: "card", ID: cardID}
}

func (m *mockCardRepo) Save(_ context.Context, _ string, card *domain.Card) error {
	m.card = card
	return m.saveErr
}

func (m *mockCardRepo) Delete(_ context.Context, _, _ string) error {
	return m.delErr
}

func (m *mockCardRepo) NextID(_ context.Context, _ string) (string, error) {
	if m.nextIDErr != nil {
		return "", m.nextIDErr
	}
	return m.nextID, nil
}

func (m *mockCardRepo) Create(_ context.Context, _ string, card *domain.Card) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	card.ID = m.nextID
	m.card = card
	return m.nextID, nil
}

func newCardRouter(cardRepo *mockCardRepo, boardRepo *mockBoardRepo) *chi.Mux {
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)
	h := handler.NewCardHandler(uc)
	r := chi.NewRouter()
	h.Register(r)
	return r
}

func TestCardHandler_List(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{cards: []domain.Card{
		{ID: "card-1", Title: "Test", List: "todo"},
	}}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/board-1/cards", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var cards []domain.Card
	if err := json.NewDecoder(w.Body).Decode(&cards); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("got %d cards, want 1", len(cards))
	}
}

func TestCardHandler_Create(t *testing.T) {
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}
	cardRepo := &mockCardRepo{nextID: "20260124-001"}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"title":"New Card","list":"todo"}`
	req := httptest.NewRequest(http.MethodPost, "/api/boards/board-1/cards", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestCardHandler_Get(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Test", List: "todo"},
	}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/board-1/cards/card-1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCardHandler_Get_NotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/board-1/cards/missing", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_Move(t *testing.T) {
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}, {ID: "done", Name: "Done"}}},
	}
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Test", List: "todo"},
	}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"list":"done","order":0}`
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/card-1/move", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestCardHandler_Move_MissingList(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"order":0}`
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/card-1/move", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_Archive(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Test", List: "todo", Archived: false},
	}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"archived":true}`
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/card-1/archive", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestCardHandler_Delete(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Test", List: "todo"},
	}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/board-1/cards/card-1", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestCardHandler_Update(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Original", List: "todo"},
	}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"title":"Updated Title","description":"New desc"}`
	req := httptest.NewRequest(http.MethodPut, "/api/boards/board-1/cards/card-1", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestCardHandler_Create_InvalidBody(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/boards/board-1/cards", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_Update_InvalidBody(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodPut, "/api/boards/board-1/cards/card-1", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_Move_InvalidBody(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/card-1/move", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_Archive_InvalidBody(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/card-1/archive", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_Create_ValidationError(t *testing.T) {
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}
	cardRepo := &mockCardRepo{nextID: "20260124-001"}
	r := newCardRouter(cardRepo, boardRepo)

	// Missing title triggers validation error
	body := `{"list":"todo"}`
	req := httptest.NewRequest(http.MethodPost, "/api/boards/board-1/cards", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d. body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

func TestCardHandler_Delete_NotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodDelete, "/api/boards/board-1/cards/missing", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_Update_NotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"title":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/boards/board-1/cards/missing", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_Move_CardNotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "done", Name: "Done"}}},
	}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"list":"done","order":0}`
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/missing/move", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_Archive_NotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	body := `{"archived":true}`
	req := httptest.NewRequest(http.MethodPatch, "/api/boards/board-1/cards/missing/archive", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_List_BoardNotFound(t *testing.T) {
	boardRepo := &mockBoardRepo{}
	cardRepo := &mockCardRepo{}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/missing/cards", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCardHandler_List_WithArchived(t *testing.T) {
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	cardRepo := &mockCardRepo{cards: []domain.Card{
		{ID: "card-1", Title: "Active", List: "todo"},
		{ID: "card-2", Title: "Archived", List: "todo", Archived: true},
	}}
	r := newCardRouter(cardRepo, boardRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/boards/board-1/cards?archived=true", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var cards []domain.Card
	if err := json.NewDecoder(w.Body).Decode(&cards); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("got %d cards, want 2", len(cards))
	}
}
