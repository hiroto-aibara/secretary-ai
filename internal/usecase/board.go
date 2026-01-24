package usecase

import (
	"context"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

type BoardUseCase struct {
	repo domain.BoardRepository
}

func NewBoardUseCase(repo domain.BoardRepository) *BoardUseCase {
	return &BoardUseCase{repo: repo}
}

func (uc *BoardUseCase) List(ctx context.Context) ([]domain.Board, error) {
	return uc.repo.List(ctx)
}

func (uc *BoardUseCase) Get(ctx context.Context, id string) (*domain.Board, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *BoardUseCase) Create(ctx context.Context, board *domain.Board) (*domain.Board, error) {
	if err := board.Validate(); err != nil {
		return nil, err
	}

	existing, err := uc.repo.Get(ctx, board.ID)
	if err == nil && existing != nil {
		return nil, &domain.ErrConflict{Resource: "board", ID: board.ID}
	}

	if err := uc.repo.Save(ctx, board); err != nil {
		return nil, err
	}
	return board, nil
}

func (uc *BoardUseCase) Update(ctx context.Context, id string, board *domain.Board) (*domain.Board, error) {
	existing, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if board.Name != "" {
		existing.Name = board.Name
	}
	if len(board.Lists) > 0 {
		existing.Lists = board.Lists
	}

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *BoardUseCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.Get(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}
