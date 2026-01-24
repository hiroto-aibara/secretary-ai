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

	card.List = toList
	card.Order = order
	card.UpdatedAt = time.Now()

	if err := uc.cardRepo.Save(ctx, boardID, card); err != nil {
		return nil, err
	}
	return card, nil
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
