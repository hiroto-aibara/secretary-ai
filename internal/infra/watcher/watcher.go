package watcher

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type event struct {
	Type    string `json:"type"`
	BoardID string `json:"board_id"`
	Time    string `json:"timestamp"`
}

type Broadcaster interface {
	BroadcastRaw(data []byte)
}

type Watcher struct {
	broadcaster Broadcaster
	basePath    string
	debounce    time.Duration
}

func New(broadcaster Broadcaster, basePath string) *Watcher {
	return &Watcher{
		broadcaster: broadcaster,
		basePath:    basePath,
		debounce:    500 * time.Millisecond,
	}
}

func (w *Watcher) Start() error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fsw.Close()

	boardsDir := filepath.Join(w.basePath, "boards")
	if err := w.addRecursive(fsw, boardsDir); err != nil {
		slog.Warn("watcher: initial watch setup", "error", err)
	}

	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	var pendingEvent *event

	for {
		select {
		case fsEvent, ok := <-fsw.Events:
			if !ok {
				return nil
			}
			slog.Debug("file changed", "path", fsEvent.Name, "op", fsEvent.Op.String())

			if fsEvent.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}

			if !strings.HasSuffix(fsEvent.Name, ".yaml") {
				if fsEvent.Op&fsnotify.Create != 0 {
					_ = w.addRecursive(fsw, fsEvent.Name)
				}
				continue
			}

			ev := w.classifyEvent(fsEvent.Name)
			if ev == nil {
				continue
			}
			pendingEvent = ev

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(w.debounce)

		case <-timer.C:
			if pendingEvent != nil {
				data, err := json.Marshal(pendingEvent)
				if err != nil {
					slog.Error("failed to marshal watcher event", "error", err)
				} else {
					w.broadcaster.BroadcastRaw(data)
				}
				pendingEvent = nil
			}

		case err, ok := <-fsw.Errors:
			if !ok {
				return nil
			}
			slog.Error("watcher error", "error", err)
		}
	}
}

func (w *Watcher) classifyEvent(path string) *event {
	rel, err := filepath.Rel(w.basePath, path)
	if err != nil {
		return nil
	}

	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 3 || parts[0] != "boards" {
		return nil
	}

	boardID := parts[1]
	eventType := "board_updated"

	if len(parts) >= 4 && parts[2] == "cards" {
		eventType = "card_updated"
	}

	return &event{
		Type:    eventType,
		BoardID: boardID,
		Time:    time.Now().Format(time.RFC3339),
	}
}

func (w *Watcher) addRecursive(fsw *fsnotify.Watcher, path string) error {
	return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return fsw.Add(p)
		}
		return nil
	})
}
