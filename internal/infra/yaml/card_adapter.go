package yaml

import (
	"context"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

// CardRepositoryAdapter adapts Store to satisfy domain.CardRepository interface.
type CardRepositoryAdapter struct {
	store *Store
}

func NewCardRepositoryAdapter(store *Store) *CardRepositoryAdapter {
	return &CardRepositoryAdapter{store: store}
}

func (a *CardRepositoryAdapter) ListByBoard(ctx context.Context, boardID string, includeArchived bool) ([]domain.Card, error) {
	return a.store.ListByBoard(ctx, boardID, includeArchived)
}

func (a *CardRepositoryAdapter) Get(ctx context.Context, boardID, cardID string) (*domain.Card, error) {
	return a.store.GetCard(ctx, boardID, cardID)
}

func (a *CardRepositoryAdapter) Save(ctx context.Context, boardID string, card *domain.Card) error {
	return a.store.SaveCard(ctx, boardID, card)
}

func (a *CardRepositoryAdapter) Delete(ctx context.Context, boardID, cardID string) error {
	return a.store.DeleteCard(ctx, boardID, cardID)
}

func (a *CardRepositoryAdapter) NextID(ctx context.Context, boardID string) (string, error) {
	return a.store.NextID(ctx, boardID)
}
