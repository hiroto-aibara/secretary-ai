package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
)

type mockCardRepo struct {
	cards       []domain.Card
	card        *domain.Card
	getErr      error
	saveErr     error
	delErr      error
	nextID      string
	nextIDErr   error
	createErr   error
	savedCard   *domain.Card
	createdCard *domain.Card
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
	m.savedCard = card
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
	m.createdCard = card
	return m.nextID, nil
}

func TestCardUseCase_List(t *testing.T) {
	cards := []domain.Card{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	cardRepo := &mockCardRepo{cards: cards}
	boardRepo := &mockBoardRepo{board: &domain.Board{ID: "board-1"}}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	got, err := uc.List(context.Background(), "board-1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d cards, want 2", len(got))
	}
}

func TestCardUseCase_List_BoardNotFound(t *testing.T) {
	cardRepo := &mockCardRepo{}
	boardRepo := &mockBoardRepo{}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.List(context.Background(), "missing", false)
	if err == nil {
		t.Error("expected error, got nil")
	}
	var notFound *domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestCardUseCase_Create(t *testing.T) {
	tests := []struct {
		name      string
		card      *domain.Card
		setup     func(*mockCardRepo, *mockBoardRepo)
		wantErr   bool
		checkType func(error) bool
	}{
		{
			name: "success",
			card: &domain.Card{Title: "New Card", List: "todo"},
			setup: func(cr *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
				cr.nextID = "20260124-001"
			},
		},
		{
			name: "invalid list",
			card: &domain.Card{Title: "New Card", List: "invalid"},
			setup: func(_ *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
			},
			wantErr: true,
			checkType: func(err error) bool {
				var ve *domain.ErrValidation
				return errors.As(err, &ve)
			},
		},
		{
			name: "missing title",
			card: &domain.Card{List: "todo"},
			setup: func(_ *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
			},
			wantErr: true,
			checkType: func(err error) bool {
				var ve *domain.ErrValidation
				return errors.As(err, &ve)
			},
		},
		{
			name: "board not found",
			card: &domain.Card{Title: "New Card", List: "todo"},
			setup: func(_ *mockCardRepo, _ *mockBoardRepo) {
			},
			wantErr: true,
			checkType: func(err error) bool {
				var nf *domain.ErrNotFound
				return errors.As(err, &nf)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			boardRepo := &mockBoardRepo{}
			tt.setup(cardRepo, boardRepo)
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			_, err := uc.Create(context.Background(), "board-1", tt.card)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.checkType != nil && !tt.checkType(err) {
					t.Errorf("unexpected error type: %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCardUseCase_Update(t *testing.T) {
	tests := []struct {
		name    string
		cardID  string
		updates *domain.Card
		setup   func(*mockCardRepo)
		wantErr bool
	}{
		{
			name:    "success",
			cardID:  "card-1",
			updates: &domain.Card{Title: "Updated"},
			setup: func(m *mockCardRepo) {
				m.card = &domain.Card{ID: "card-1", Title: "Original", List: "todo"}
			},
		},
		{
			name:    "not found",
			cardID:  "missing",
			updates: &domain.Card{Title: "Updated"},
			setup:   func(_ *mockCardRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			tt.setup(cardRepo)
			boardRepo := &mockBoardRepo{}
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			got, err := uc.Update(context.Background(), "board-1", tt.cardID, tt.updates)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Title != "Updated" {
				t.Errorf("Title = %s, want Updated", got.Title)
			}
		})
	}
}

func TestCardUseCase_Move(t *testing.T) {
	tests := []struct {
		name      string
		cardID    string
		toList    string
		order     int
		setup     func(*mockCardRepo, *mockBoardRepo)
		wantErr   bool
		checkType func(error) bool
	}{
		{
			name:   "success",
			cardID: "card-1",
			toList: "done",
			order:  0,
			setup: func(cr *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}, {ID: "done", Name: "Done"}}}
				cr.card = &domain.Card{ID: "card-1", List: "todo", Title: "Test"}
			},
		},
		{
			name:   "invalid list",
			cardID: "card-1",
			toList: "invalid",
			order:  0,
			setup: func(cr *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
				cr.card = &domain.Card{ID: "card-1", List: "todo", Title: "Test"}
			},
			wantErr: true,
			checkType: func(err error) bool {
				var ve *domain.ErrValidation
				return errors.As(err, &ve)
			},
		},
		{
			name:   "card not found",
			cardID: "missing",
			toList: "done",
			order:  0,
			setup: func(_ *mockCardRepo, br *mockBoardRepo) {
				br.board = &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "done", Name: "Done"}}}
			},
			wantErr: true,
			checkType: func(err error) bool {
				var nf *domain.ErrNotFound
				return errors.As(err, &nf)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			boardRepo := &mockBoardRepo{}
			tt.setup(cardRepo, boardRepo)
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			got, err := uc.Move(context.Background(), "board-1", tt.cardID, tt.toList, tt.order)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.checkType != nil && !tt.checkType(err) {
					t.Errorf("unexpected error type: %T", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.List != tt.toList {
				t.Errorf("List = %s, want %s", got.List, tt.toList)
			}
		})
	}
}

func TestCardUseCase_Archive(t *testing.T) {
	tests := []struct {
		name         string
		cardID       string
		archived     bool
		setup        func(*mockCardRepo)
		wantArchived bool
		wantErr      bool
	}{
		{
			name:     "archive active card",
			cardID:   "card-1",
			archived: true,
			setup: func(m *mockCardRepo) {
				m.card = &domain.Card{ID: "card-1", Archived: false, Title: "Test", List: "todo"}
			},
			wantArchived: true,
		},
		{
			name:     "restore archived card",
			cardID:   "card-1",
			archived: false,
			setup: func(m *mockCardRepo) {
				m.card = &domain.Card{ID: "card-1", Archived: true, Title: "Test", List: "todo"}
			},
			wantArchived: false,
		},
		{
			name:     "not found",
			cardID:   "missing",
			archived: true,
			setup:    func(_ *mockCardRepo) {},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			tt.setup(cardRepo)
			boardRepo := &mockBoardRepo{}
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			got, err := uc.Archive(context.Background(), "board-1", tt.cardID, tt.archived)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Archived != tt.wantArchived {
				t.Errorf("Archived = %v, want %v", got.Archived, tt.wantArchived)
			}
		})
	}
}

func TestCardUseCase_Get(t *testing.T) {
	tests := []struct {
		name    string
		cardID  string
		setup   func(*mockCardRepo)
		wantErr bool
	}{
		{
			name:   "success",
			cardID: "card-1",
			setup: func(m *mockCardRepo) {
				m.card = &domain.Card{ID: "card-1", Title: "Test", List: "todo"}
			},
		},
		{
			name:    "not found",
			cardID:  "missing",
			setup:   func(_ *mockCardRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			tt.setup(cardRepo)
			boardRepo := &mockBoardRepo{}
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			got, err := uc.Get(context.Background(), "board-1", tt.cardID)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.cardID {
				t.Errorf("ID = %s, want %s", got.ID, tt.cardID)
			}
		})
	}
}

func TestCardUseCase_Create_Error(t *testing.T) {
	cardRepo := &mockCardRepo{
		createErr: errors.New("create failed"),
	}
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
	}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.Create(context.Background(), "board-1", &domain.Card{Title: "New", List: "todo"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCardUseCase_Update_SaveError(t *testing.T) {
	cardRepo := &mockCardRepo{
		card:    &domain.Card{ID: "card-1", Title: "Original", List: "todo"},
		saveErr: errors.New("save failed"),
	}
	boardRepo := &mockBoardRepo{}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.Update(context.Background(), "board-1", "card-1", &domain.Card{Title: "Updated"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCardUseCase_Move_BoardNotFound(t *testing.T) {
	cardRepo := &mockCardRepo{}
	boardRepo := &mockBoardRepo{}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.Move(context.Background(), "missing", "card-1", "done", 0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCardUseCase_Move_SaveError(t *testing.T) {
	cardRepo := &mockCardRepo{
		card:    &domain.Card{ID: "card-1", Title: "Test", List: "todo"},
		saveErr: errors.New("save failed"),
	}
	boardRepo := &mockBoardRepo{
		board: &domain.Board{ID: "board-1", Lists: []domain.List{{ID: "todo", Name: "Todo"}, {ID: "done", Name: "Done"}}},
	}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.Move(context.Background(), "board-1", "card-1", "done", 0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCardUseCase_Archive_SaveError(t *testing.T) {
	cardRepo := &mockCardRepo{
		card:    &domain.Card{ID: "card-1", Title: "Test", List: "todo", Archived: false},
		saveErr: errors.New("save failed"),
	}
	boardRepo := &mockBoardRepo{}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	_, err := uc.Archive(context.Background(), "board-1", "card-1", true)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCardUseCase_Update_PartialFields(t *testing.T) {
	cardRepo := &mockCardRepo{
		card: &domain.Card{ID: "card-1", Title: "Original", Description: "Desc", List: "todo", Labels: []string{"bug"}},
	}
	boardRepo := &mockBoardRepo{}
	uc := usecase.NewCardUseCase(cardRepo, boardRepo)

	got, err := uc.Update(context.Background(), "board-1", "card-1", &domain.Card{Labels: []string{"feature"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "Original" {
		t.Errorf("Title = %s, want Original", got.Title)
	}
	if len(got.Labels) != 1 || got.Labels[0] != "feature" {
		t.Errorf("Labels = %v, want [feature]", got.Labels)
	}
}

func TestCardUseCase_Delete(t *testing.T) {
	tests := []struct {
		name    string
		cardID  string
		setup   func(*mockCardRepo)
		wantErr bool
	}{
		{
			name:   "success",
			cardID: "card-1",
			setup: func(m *mockCardRepo) {
				m.card = &domain.Card{ID: "card-1", Title: "Test", List: "todo"}
			},
		},
		{
			name:    "not found",
			cardID:  "missing",
			setup:   func(_ *mockCardRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardRepo := &mockCardRepo{}
			tt.setup(cardRepo)
			boardRepo := &mockBoardRepo{}
			uc := usecase.NewCardUseCase(cardRepo, boardRepo)

			err := uc.Delete(context.Background(), "board-1", tt.cardID)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
