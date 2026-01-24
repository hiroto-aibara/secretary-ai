package domain_test

import (
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

func TestErrNotFound_Error(t *testing.T) {
	err := &domain.ErrNotFound{Resource: "card", ID: "123"}
	want := "card 123 not found"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrValidation_Error(t *testing.T) {
	err := &domain.ErrValidation{Field: "title", Message: "is required"}
	want := "validation error: title is required"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrConflict_Error(t *testing.T) {
	err := &domain.ErrConflict{Resource: "board", ID: "test"}
	want := "board test already exists"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
