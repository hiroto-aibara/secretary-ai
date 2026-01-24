package yaml

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/hiroto-aibara/secretary-ai/internal/domain"
)

type Store struct {
	basePath string
	mu       sync.RWMutex
}

func NewStore(basePath string) *Store {
	return &Store{basePath: basePath}
}

func (s *Store) boardDir(id string) string {
	return filepath.Join(s.basePath, "boards", id)
}

func (s *Store) boardFile(id string) string {
	return filepath.Join(s.boardDir(id), "board.yaml")
}

func (s *Store) cardsDir(boardID string) string {
	return filepath.Join(s.boardDir(boardID), "cards")
}

func (s *Store) cardFile(boardID, cardID string) string {
	return filepath.Join(s.cardsDir(boardID), cardID+".yaml")
}

// BoardRepository implementation

func (s *Store) List(_ context.Context) ([]domain.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	boardsDir := filepath.Join(s.basePath, "boards")
	entries, err := os.ReadDir(boardsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Board{}, nil
		}
		return nil, fmt.Errorf("read boards dir: %w", err)
	}

	var boards []domain.Board
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		board, err := s.readBoard(entry.Name())
		if err != nil {
			continue
		}
		boards = append(boards, *board)
	}
	return boards, nil
}

func (s *Store) Get(_ context.Context, id string) (*domain.Board, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.readBoard(id)
}

func (s *Store) Save(_ context.Context, board *domain.Board) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := s.boardDir(board.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create board dir: %w", err)
	}
	if err := os.MkdirAll(s.cardsDir(board.ID), 0o755); err != nil {
		return fmt.Errorf("create cards dir: %w", err)
	}

	data, err := yamlv3.Marshal(board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}

	if err := os.WriteFile(s.boardFile(board.ID), data, 0o644); err != nil {
		return fmt.Errorf("write board file: %w", err)
	}
	return nil
}

func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := s.boardDir(id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return &domain.ErrNotFound{Resource: "board", ID: id}
	}
	return os.RemoveAll(dir)
}

func (s *Store) readBoard(id string) (*domain.Board, error) {
	data, err := os.ReadFile(s.boardFile(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &domain.ErrNotFound{Resource: "board", ID: id}
		}
		return nil, fmt.Errorf("read board file: %w", err)
	}

	var board domain.Board
	if err := yamlv3.Unmarshal(data, &board); err != nil {
		return nil, fmt.Errorf("unmarshal board: %w", err)
	}
	return &board, nil
}

// CardRepository implementation

func (s *Store) ListByBoard(_ context.Context, boardID string, includeArchived bool) ([]domain.Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dir := s.cardsDir(boardID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Card{}, nil
		}
		return nil, fmt.Errorf("read cards dir: %w", err)
	}

	var cards []domain.Card
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		card, err := s.readCard(boardID, strings.TrimSuffix(entry.Name(), ".yaml"))
		if err != nil {
			continue
		}
		if !includeArchived && card.Archived {
			continue
		}
		cards = append(cards, *card)
	}

	sort.Slice(cards, func(i, j int) bool {
		if cards[i].List != cards[j].List {
			return cards[i].List < cards[j].List
		}
		return cards[i].Order < cards[j].Order
	})

	return cards, nil
}

func (s *Store) GetCard(_ context.Context, boardID, cardID string) (*domain.Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.readCard(boardID, cardID)
}

func (s *Store) SaveCard(_ context.Context, boardID string, card *domain.Card) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := s.cardsDir(boardID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create cards dir: %w", err)
	}

	data, err := yamlv3.Marshal(card)
	if err != nil {
		return fmt.Errorf("marshal card: %w", err)
	}

	if err := os.WriteFile(s.cardFile(boardID, card.ID), data, 0o644); err != nil {
		return fmt.Errorf("write card file: %w", err)
	}
	return nil
}

func (s *Store) DeleteCard(_ context.Context, boardID, cardID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.cardFile(boardID, cardID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &domain.ErrNotFound{Resource: "card", ID: cardID}
	}
	return os.Remove(path)
}

func (s *Store) NextID(_ context.Context, boardID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nextIDLocked(boardID)
}

func (s *Store) nextIDLocked(boardID string) (string, error) {
	today := time.Now().Format("20060102")
	dir := s.cardsDir(boardID)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return today + "-001", nil
		}
		return "", fmt.Errorf("read cards dir: %w", err)
	}

	maxSeq := 0
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		if strings.HasPrefix(name, today+"-") {
			seqStr := strings.TrimPrefix(name, today+"-")
			seq := 0
			for _, ch := range seqStr {
				if ch >= '0' && ch <= '9' {
					seq = seq*10 + int(ch-'0')
				}
			}
			if seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	return fmt.Sprintf("%s-%03d", today, maxSeq+1), nil
}

func (s *Store) CreateCard(_ context.Context, boardID string, card *domain.Card) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := s.nextIDLocked(boardID)
	if err != nil {
		return "", err
	}
	card.ID = id

	dir := s.cardsDir(boardID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create cards dir: %w", err)
	}

	data, err := yamlv3.Marshal(card)
	if err != nil {
		return "", fmt.Errorf("marshal card: %w", err)
	}

	if err := os.WriteFile(s.cardFile(boardID, card.ID), data, 0o644); err != nil {
		return "", fmt.Errorf("write card file: %w", err)
	}
	return id, nil
}

func (s *Store) readCard(boardID, cardID string) (*domain.Card, error) {
	data, err := os.ReadFile(s.cardFile(boardID, cardID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &domain.ErrNotFound{Resource: "card", ID: cardID}
		}
		return nil, fmt.Errorf("read card file: %w", err)
	}

	var card domain.Card
	if err := yamlv3.Unmarshal(data, &card); err != nil {
		return nil, fmt.Errorf("unmarshal card: %w", err)
	}
	return &card, nil
}

// BasePath returns the store's base path for external use (e.g., file watcher).
func (s *Store) BasePath() string {
	return s.basePath
}
