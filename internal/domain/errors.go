package domain

import "fmt"

type ErrNotFound struct {
	Resource string
	ID       string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}

type ErrValidation struct {
	Field   string
	Message string
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

type ErrConflict struct {
	Resource string
	ID       string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("%s %s already exists", e.Resource, e.ID)
}
