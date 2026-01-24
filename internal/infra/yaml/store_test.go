package yaml_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
	yamlstore "github.com/hiroto-aibara/secretary-ai/internal/infra/yaml"
)

func setupStore(t *testing.T) *yamlstore.Store {
	t.Helper()
	dir := t.TempDir()
	return yamlstore.NewStore(dir)
}

func TestStore_Board_CRUD(t *testing.T) {
	store := setupStore(t)
	ctx := context.Background()

	// List empty
	boards, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(boards) != 0 {
		t.Errorf("got %d boards, want 0", len(boards))
	}

	// Save
	board := &domain.Board{
		ID:    "test-board",
		Name:  "Test Board",
		Lists: []domain.List{{ID: "todo", Name: "Todo"}, {ID: "done", Name: "Done"}},
	}
	if err := store.Save(ctx, board); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Get
	got, err := store.Get(ctx, "test-board")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Test Board" {
		t.Errorf("Name = %s, want Test Board", got.Name)
	}
	if len(got.Lists) != 2 {
		t.Errorf("Lists count = %d, want 2", len(got.Lists))
	}

	// List
	boards, err = store.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(boards) != 1 {
		t.Errorf("got %d boards, want 1", len(boards))
	}

	// Delete
	if err := store.Delete(ctx, "test-board"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Get after delete
	_, err = store.Get(ctx, "test-board")
	var notFound *domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestStore_Board_GetNotFound(t *testing.T) {
	store := setupStore(t)
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	var notFound *domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStore_Card_CRUD(t *testing.T) {
	store := setupStore(t)
	adapter := yamlstore.NewCardRepositoryAdapter(store)
	ctx := context.Background()

	// Create board first
	board := &domain.Board{
		ID:    "board-1",
		Name:  "Board",
		Lists: []domain.List{{ID: "todo", Name: "Todo"}},
	}
	if err := store.Save(ctx, board); err != nil {
		t.Fatalf("Save board: %v", err)
	}

	// List empty
	cards, err := adapter.ListByBoard(ctx, "board-1", false)
	if err != nil {
		t.Fatalf("ListByBoard: %v", err)
	}
	if len(cards) != 0 {
		t.Errorf("got %d cards, want 0", len(cards))
	}

	// Save card
	card := &domain.Card{
		ID:    "20260124-001",
		Title: "Test Card",
		List:  "todo",
		Order: 0,
	}
	if err := adapter.Save(ctx, "board-1", card); err != nil {
		t.Fatalf("Save card: %v", err)
	}

	// Get card
	got, err := adapter.Get(ctx, "board-1", "20260124-001")
	if err != nil {
		t.Fatalf("Get card: %v", err)
	}
	if got.Title != "Test Card" {
		t.Errorf("Title = %s, want Test Card", got.Title)
	}

	// List cards
	cards, err = adapter.ListByBoard(ctx, "board-1", false)
	if err != nil {
		t.Fatalf("ListByBoard: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("got %d cards, want 1", len(cards))
	}

	// Delete card
	if err := adapter.Delete(ctx, "board-1", "20260124-001"); err != nil {
		t.Fatalf("Delete card: %v", err)
	}

	// Get after delete
	_, err = adapter.Get(ctx, "board-1", "20260124-001")
	var notFound *domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestStore_Card_ArchiveFilter(t *testing.T) {
	store := setupStore(t)
	adapter := yamlstore.NewCardRepositoryAdapter(store)
	ctx := context.Background()

	board := &domain.Board{ID: "board-1", Name: "Board", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
	if err := store.Save(ctx, board); err != nil {
		t.Fatalf("Save board: %v", err)
	}

	// Save active and archived cards
	if err := adapter.Save(ctx, "board-1", &domain.Card{ID: "active", Title: "Active", List: "todo", Archived: false}); err != nil {
		t.Fatal(err)
	}
	if err := adapter.Save(ctx, "board-1", &domain.Card{ID: "archived", Title: "Archived", List: "todo", Archived: true}); err != nil {
		t.Fatal(err)
	}

	// Without archived
	cards, err := adapter.ListByBoard(ctx, "board-1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Errorf("without archived: got %d cards, want 1", len(cards))
	}

	// With archived
	cards, err = adapter.ListByBoard(ctx, "board-1", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 2 {
		t.Errorf("with archived: got %d cards, want 2", len(cards))
	}
}

func TestStore_NextID(t *testing.T) {
	store := setupStore(t)
	adapter := yamlstore.NewCardRepositoryAdapter(store)
	ctx := context.Background()

	board := &domain.Board{ID: "board-1", Name: "Board", Lists: []domain.List{{ID: "todo", Name: "Todo"}}}
	if err := store.Save(ctx, board); err != nil {
		t.Fatal(err)
	}

	// First ID
	id1, err := adapter.NextID(ctx, "board-1")
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if len(id1) != 12 { // YYYYMMDD-NNN
		t.Errorf("ID length = %d, want 12. ID = %s", len(id1), id1)
	}

	// Save a card and get next ID
	if err := adapter.Save(ctx, "board-1", &domain.Card{ID: id1, Title: "First", List: "todo"}); err != nil {
		t.Fatal(err)
	}

	id2, err := adapter.NextID(ctx, "board-1")
	if err != nil {
		t.Fatalf("NextID: %v", err)
	}
	if id2 == id1 {
		t.Errorf("second ID should differ from first: %s", id2)
	}
}

func TestStore_Card_GetNotFound(t *testing.T) {
	store := setupStore(t)
	adapter := yamlstore.NewCardRepositoryAdapter(store)
	ctx := context.Background()

	_, err := adapter.Get(ctx, "board-1", "nonexistent")
	var notFound *domain.ErrNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
