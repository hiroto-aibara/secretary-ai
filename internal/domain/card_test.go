package domain_test

import (
	"errors"
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

func TestCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    domain.Card
		wantErr bool
		field   string
	}{
		{
			name:    "valid card",
			card:    domain.Card{Title: "Test", List: "todo"},
			wantErr: false,
		},
		{
			name:    "missing title",
			card:    domain.Card{List: "todo"},
			wantErr: true,
			field:   "title",
		},
		{
			name:    "missing list",
			card:    domain.Card{Title: "Test"},
			wantErr: true,
			field:   "list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
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
