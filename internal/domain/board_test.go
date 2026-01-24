package domain_test

import (
	"errors"
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

func TestBoard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		board   domain.Board
		wantErr bool
		field   string
	}{
		{
			name:    "valid board",
			board:   domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			wantErr: false,
		},
		{
			name:    "missing id",
			board:   domain.Board{Name: "Test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			wantErr: true,
			field:   "id",
		},
		{
			name:    "missing name",
			board:   domain.Board{ID: "test", Lists: []domain.List{{ID: "todo", Name: "Todo"}}},
			wantErr: true,
			field:   "name",
		},
		{
			name:    "empty lists",
			board:   domain.Board{ID: "test", Name: "Test", Lists: []domain.List{}},
			wantErr: true,
			field:   "lists",
		},
		{
			name:    "list missing id",
			board:   domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{Name: "Todo"}}},
			wantErr: true,
			field:   "lists.id",
		},
		{
			name:    "list missing name",
			board:   domain.Board{ID: "test", Name: "Test", Lists: []domain.List{{ID: "todo"}}},
			wantErr: true,
			field:   "lists.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.board.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				var ve *domain.ErrValidation
				if !errors.As(err, &ve) {
					t.Errorf("expected ErrValidation, got %T", err)
					return
				}
				if ve.Field != tt.field {
					t.Errorf("field = %s, want %s", ve.Field, tt.field)
				}
			}
		})
	}
}

func TestBoard_HasList(t *testing.T) {
	board := domain.Board{
		Lists: []domain.List{
			{ID: "todo", Name: "Todo"},
			{ID: "done", Name: "Done"},
		},
	}

	tests := []struct {
		name   string
		listID string
		want   bool
	}{
		{"existing list", "todo", true},
		{"another existing list", "done", true},
		{"non-existing list", "in-progress", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := board.HasList(tt.listID); got != tt.want {
				t.Errorf("HasList(%s) = %v, want %v", tt.listID, got, tt.want)
			}
		})
	}
}
