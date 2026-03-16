package usecase

import (
	"context"
	"time"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

type CardUseCase struct {
	cardRepo  domain.CardRepository
	boardRepo domain.BoardRepository
}

func NewCardUseCase(cardRepo domain.CardRepository, boardRepo domain.BoardRepository) *CardUseCase {
	return &CardUseCase{cardRepo: cardRepo, boardRepo: boardRepo}
}

func (uc *CardUseCase) List(ctx context.Context, boardID string, includeArchived bool) ([]domain.Card, error) {
	if _, err := uc.boardRepo.Get(ctx, boardID); err != nil {
		return nil, err
	}
	return uc.cardRepo.ListByBoard(ctx, boardID, includeArchived)
}

func (uc *CardUseCase) Get(ctx context.Context, boardID, cardID string) (*domain.Card, error) {
	return uc.cardRepo.Get(ctx, boardID, cardID)
}

func (uc *CardUseCase) Create(ctx context.Context, boardID string, card *domain.Card) (*domain.Card, error) {
	board, err := uc.boardRepo.Get(ctx, boardID)
	if err != nil {
		return nil, err
	}

	if err := card.Validate(); err != nil {
		return nil, err
	}

	if !board.HasList(card.List) {
		return nil, &domain.ErrValidation{
			Field:   "list",
			Message: "list '" + card.List + "' does not exist in board",
		}
	}

	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now
	card.Archived = false

	id, err := uc.cardRepo.Create(ctx, boardID, card)
	if err != nil {
		return nil, err
	}
	card.ID = id
	return card, nil
}

func (uc *CardUseCase) Update(ctx context.Context, boardID, cardID string, updates *domain.Card) (*domain.Card, error) {
	existing, err := uc.cardRepo.Get(ctx, boardID, cardID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Labels != nil {
		existing.Labels = updates.Labels
	}
	if updates.Todos != nil {
		existing.Todos = updates.Todos
	}

	existing.UpdatedAt = time.Now()

	if err := uc.cardRepo.Save(ctx, boardID, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *CardUseCase) Delete(ctx context.Context, boardID, cardID string) error {
	if _, err := uc.cardRepo.Get(ctx, boardID, cardID); err != nil {
		return err
	}
	return uc.cardRepo.Delete(ctx, boardID, cardID)
}

func (uc *CardUseCase) Move(ctx context.Context, boardID, cardID, toList string, order int) (*domain.Card, error) {
	board, err := uc.boardRepo.Get(ctx, boardID)
	if err != nil {
		return nil, err
	}

	if !board.HasList(toList) {
		return nil, &domain.ErrValidation{
			Field:   "list",
			Message: "list '" + toList + "' does not exist in board",
		}
	}

	card, err := uc.cardRepo.Get(ctx, boardID, cardID)
	if err != nil {
		return nil, err
	}

	fromList := card.List
	card.List = toList
	card.Order = order
	card.UpdatedAt = time.Now()

	if err := uc.cardRepo.Save(ctx, boardID, card); err != nil {
		return nil, err
	}

	if err := uc.reorderList(ctx, boardID, toList, cardID, order); err != nil {
		return nil, err
	}

	if fromList != toList {
		if err := uc.reorderList(ctx, boardID, fromList, "", -1); err != nil {
			return nil, err
		}
	}

	return card, nil
}

func (uc *CardUseCase) reorderList(ctx context.Context, boardID, listID, movedCardID string, targetOrder int) error {
	allCards, err := uc.cardRepo.ListByBoard(ctx, boardID, false)
	if err != nil {
		return err
	}

	var listCards []domain.Card
	for _, c := range allCards {
		if c.List == listID {
			listCards = append(listCards, c)
		}
	}

	if movedCardID != "" {
		var without []domain.Card
		var moved *domain.Card
		for i, c := range listCards {
			if c.ID == movedCardID {
				moved = &listCards[i]
			} else {
				without = append(without, c)
			}
		}
		if moved != nil {
			idx := targetOrder
			if idx > len(without) {
				idx = len(without)
			}
			if idx < 0 {
				idx = 0
			}
			result := make([]domain.Card, 0, len(without)+1)
			result = append(result, without[:idx]...)
			result = append(result, *moved)
			result = append(result, without[idx:]...)
			listCards = result
		}
	}

	now := time.Now()
	for i, c := range listCards {
		if c.Order != i {
			c.Order = i
			c.UpdatedAt = now
			if err := uc.cardRepo.Save(ctx, boardID, &c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *CardUseCase) Archive(ctx context.Context, boardID, cardID string, archived bool) (*domain.Card, error) {
	card, err := uc.cardRepo.Get(ctx, boardID, cardID)
	if err != nil {
		return nil, err
	}

	card.Archived = archived
	card.UpdatedAt = time.Now()

	if err := uc.cardRepo.Save(ctx, boardID, card); err != nil {
		return nil, err
	}
	return card, nil
}
