package domain

import "context"

type BoardRepository interface {
	List(ctx context.Context) ([]Board, error)
	Get(ctx context.Context, id string) (*Board, error)
	Save(ctx context.Context, board *Board) error
	Delete(ctx context.Context, id string) error
}

type CardRepository interface {
	ListByBoard(ctx context.Context, boardID string, includeArchived bool) ([]Card, error)
	Get(ctx context.Context, boardID, cardID string) (*Card, error)
	Save(ctx context.Context, boardID string, card *Card) error
	Delete(ctx context.Context, boardID, cardID string) error
	NextID(ctx context.Context, boardID string) (string, error)
	Create(ctx context.Context, boardID string, card *Card) (string, error)
}
