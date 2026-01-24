package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
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

func TestBoardUseCase_List(t *testing.T) {
	boards := []domain.Board{{ID: "a", Name: "A"}, {ID: "b", Name: "B"}}
	repo := &mockBoardRepo{boards: boards}
	uc := usecase.NewBoardUseCase(repo)

	got, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d boards, want 2", len(got))
	}
}

func TestBoardUseCase_Get(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		board   *domain.Board
		wantErr bool
	}{
		{
			name:  "found",
			id:    "test",
			board: &domain.Board{ID: "test", Name: "Test"},
		},
		{
			name:    "not found",
			id:      "missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBoardRepo{board: tt.board}
			uc := usecase.NewBoardUseCase(repo)

			got, err := uc.Get(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.id {
				t.Errorf("ID = %s, want %s", got.ID, tt.id)
			}
		})
	}
}

func TestBoardUseCase_Create(t *testing.T) {
	tests := []struct {
		name      string
		board     *domain.Board
		setup     func(*mockBoardRepo)
		wantErr   bool
		checkType func(error) bool
	}{
		{
			name:  "success",
			board: &domain.Board{ID: "new", Name: "New", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			setup: func(m *mockBoardRepo) {
				m.getErr = &domain.ErrNotFound{Resource: "board", ID: "new"}
			},
		},
		{
			name:    "validation error - missing name",
			board:   &domain.Board{ID: "new", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			setup:   func(_ *mockBoardRepo) {},
			wantErr: true,
			checkType: func(err error) bool {
				var ve *domain.ErrValidation
				return errors.As(err, &ve)
			},
		},
		{
			name:  "conflict",
			board: &domain.Board{ID: "existing", Name: "Existing", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			setup: func(m *mockBoardRepo) {
				m.board = &domain.Board{ID: "existing", Name: "Existing"}
			},
			wantErr: true,
			checkType: func(err error) bool {
				var ce *domain.ErrConflict
				return errors.As(err, &ce)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBoardRepo{}
			tt.setup(repo)
			uc := usecase.NewBoardUseCase(repo)

			_, err := uc.Create(context.Background(), tt.board)
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

func TestBoardUseCase_Update(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		updates *domain.Board
		setup   func(*mockBoardRepo)
		wantErr bool
	}{
		{
			name:    "success - update name",
			id:      "test",
			updates: &domain.Board{Name: "Updated"},
			setup: func(m *mockBoardRepo) {
				m.board = &domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
			},
		},
		{
			name:    "not found",
			id:      "missing",
			updates: &domain.Board{Name: "Updated"},
			setup:   func(_ *mockBoardRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBoardRepo{}
			tt.setup(repo)
			uc := usecase.NewBoardUseCase(repo)

			got, err := uc.Update(context.Background(), tt.id, tt.updates)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != "Updated" {
				t.Errorf("Name = %s, want Updated", got.Name)
			}
		})
	}
}

func TestBoardUseCase_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(*mockBoardRepo)
		wantErr bool
	}{
		{
			name: "success",
			id:   "test",
			setup: func(m *mockBoardRepo) {
				m.board = &domain.Board{ID: "test", Name: "Test"}
			},
		},
		{
			name:    "not found",
			id:      "missing",
			setup:   func(_ *mockBoardRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBoardRepo{}
			tt.setup(repo)
			uc := usecase.NewBoardUseCase(repo)

			err := uc.Delete(context.Background(), tt.id)
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
