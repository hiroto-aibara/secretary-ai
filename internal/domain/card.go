package domain

import "time"

// TodoItem represents a single todo item within a card
type TodoItem struct {
	ID        string `json:"id" yaml:"id"`
	Text      string `json:"text" yaml:"text"`
	Completed bool   `json:"completed" yaml:"completed"`
}

type Card struct {
	ID          string     `json:"id" yaml:"id"`
	Title       string     `json:"title" yaml:"title"`
	List        string     `json:"list" yaml:"list"`
	Order       int        `json:"order" yaml:"order"`
	Description string     `json:"description" yaml:"description"`
	Labels      []string   `json:"labels" yaml:"labels"`
	Todos       []TodoItem `json:"todos" yaml:"todos"`
	Archived    bool       `json:"archived" yaml:"archived"`
	CreatedAt   time.Time  `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" yaml:"updated_at"`
}

func (c *Card) Validate() error {
	if c.Title == "" {
		return &ErrValidation{Field: "title", Message: "is required"}
	}
	if c.List == "" {
		return &ErrValidation{Field: "list", Message: "is required"}
	}
	return nil
}
