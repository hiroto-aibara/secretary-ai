package domain

type List struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type Board struct {
	ID    string `json:"id" yaml:"id"`
	Name  string `json:"name" yaml:"name"`
	Lists []List `json:"lists" yaml:"lists"`
}

func (b *Board) Validate() error {
	if b.ID == "" {
		return &ErrValidation{Field: "id", Message: "is required"}
	}
	if b.Name == "" {
		return &ErrValidation{Field: "name", Message: "is required"}
	}
	if len(b.Lists) == 0 {
		return &ErrValidation{Field: "lists", Message: "must have at least one list"}
	}
	for _, l := range b.Lists {
		if l.ID == "" {
			return &ErrValidation{Field: "lists.id", Message: "is required"}
		}
		if l.Name == "" {
			return &ErrValidation{Field: "lists.name", Message: "is required"}
		}
	}
	return nil
}

func (b *Board) HasList(listID string) bool {
	for _, l := range b.Lists {
		if l.ID == listID {
			return true
		}
	}
	return false
}
