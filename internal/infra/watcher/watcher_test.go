package watcher_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/hiroto-aibara/secretary-ai/internal/infra/watcher"
)

type mockBroadcaster struct {
	mu       sync.Mutex
	messages [][]byte
}

func (m *mockBroadcaster) BroadcastRaw(data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, data)
}

func (m *mockBroadcaster) getMessages() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([][]byte, len(m.messages))
	copy(cp, m.messages)
	return cp
}

func TestWatcher_Start_BoardUpdated(t *testing.T) {
	tmpDir := t.TempDir()
	boardsDir := filepath.Join(tmpDir, "boards", "test-board")
	if err := os.MkdirAll(boardsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create initial board.yaml
	if err := os.WriteFile(filepath.Join(boardsDir, "board.yaml"), []byte("name: Test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	bc := &mockBroadcaster{}
	w := watcher.New(bc, tmpDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start watcher in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(ctx)
	}()

	// Wait for watcher to be ready
	time.Sleep(200 * time.Millisecond)

	// Modify the board file
	if err := os.WriteFile(filepath.Join(boardsDir, "board.yaml"), []byte("name: Updated"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Wait for debounce + processing
	time.Sleep(1 * time.Second)

	msgs := bc.getMessages()
	if len(msgs) == 0 {
		t.Fatal("expected at least one broadcast message")
	}

	var ev struct {
		Type    string `json:"type"`
		BoardID string `json:"board_id"`
	}
	if err := json.Unmarshal(msgs[0], &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.Type != "board_updated" {
		t.Errorf("type = %s, want board_updated", ev.Type)
	}
	if ev.BoardID != "test-board" {
		t.Errorf("board_id = %s, want test-board", ev.BoardID)
	}
}

func TestWatcher_Start_CardUpdated(t *testing.T) {
	tmpDir := t.TempDir()
	cardsDir := filepath.Join(tmpDir, "boards", "test-board", "cards")
	if err := os.MkdirAll(cardsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	bc := &mockBroadcaster{}
	w := watcher.New(bc, tmpDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	// Create a card file
	if err := os.WriteFile(filepath.Join(cardsDir, "20260124-001.yaml"), []byte("title: Test Card"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	time.Sleep(1 * time.Second)

	msgs := bc.getMessages()
	if len(msgs) == 0 {
		t.Fatal("expected at least one broadcast message")
	}

	var ev struct {
		Type    string `json:"type"`
		BoardID string `json:"board_id"`
	}
	if err := json.Unmarshal(msgs[0], &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.Type != "card_updated" {
		t.Errorf("type = %s, want card_updated", ev.Type)
	}
	if ev.BoardID != "test-board" {
		t.Errorf("board_id = %s, want test-board", ev.BoardID)
	}
}

func TestWatcher_Start_NonYamlIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	boardsDir := filepath.Join(tmpDir, "boards", "test-board")
	if err := os.MkdirAll(boardsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	bc := &mockBroadcaster{}
	w := watcher.New(bc, tmpDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(ctx)
	}()

	time.Sleep(200 * time.Millisecond)

	// Create a non-yaml file
	if err := os.WriteFile(filepath.Join(boardsDir, "notes.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	time.Sleep(800 * time.Millisecond)

	msgs := bc.getMessages()
	if len(msgs) != 0 {
		t.Errorf("expected no broadcast for non-yaml file, got %d messages", len(msgs))
	}
}

func TestWatcher_New(t *testing.T) {
	bc := &mockBroadcaster{}
	w := watcher.New(bc, "/tmp/test")
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
