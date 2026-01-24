package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/hiroto-aibara/secretary-ai/internal/handler"
)

func TestHub_BroadcastRaw_NoClients(t *testing.T) {
	hub := handler.NewHub()
	// Should not panic with no clients
	hub.BroadcastRaw([]byte(`{"type":"test"}`))
}

func TestWSHandler_FullFlow(t *testing.T) {
	hub := handler.NewHub()
	wsH := handler.NewWSHandler(hub)

	r := chi.NewRouter()
	wsH.Register(r)

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	// Give the server time to register the client
	time.Sleep(50 * time.Millisecond)

	// Broadcast a message
	msg := `{"type":"card_updated","board_id":"test","timestamp":"2026-01-24T10:00:00Z"}`
	hub.BroadcastRaw([]byte(msg))

	// Read the broadcasted message
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	if string(data) != msg {
		t.Errorf("got %s, want %s", string(data), msg)
	}

	// Close client connection
	conn.Close()

	// Give time for disconnect to process
	time.Sleep(50 * time.Millisecond)

	// Broadcast again should not panic (client removed)
	hub.BroadcastRaw([]byte(`{"type":"test2"}`))
}

func TestWSHandler_MultipleClients(t *testing.T) {
	hub := handler.NewHub()
	wsH := handler.NewWSHandler(hub)

	r := chi.NewRouter()
	wsH.Register(r)

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"

	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial conn1: %v", err)
	}
	defer conn1.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial conn2: %v", err)
	}
	defer conn2.Close()

	time.Sleep(50 * time.Millisecond)

	msg := `{"type":"board_updated","board_id":"b1"}`
	hub.BroadcastRaw([]byte(msg))

	// Both clients should receive the message
	_ = conn1.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data1, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("read conn1: %v", err)
	}
	if string(data1) != msg {
		t.Errorf("conn1 got %s, want %s", string(data1), msg)
	}

	_ = conn2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data2, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("read conn2: %v", err)
	}
	if string(data2) != msg {
		t.Errorf("conn2 got %s, want %s", string(data2), msg)
	}
}

func TestWSHandler_UpgradeFail(t *testing.T) {
	hub := handler.NewHub()
	wsH := handler.NewWSHandler(hub)

	r := chi.NewRouter()
	wsH.Register(r)

	// Send a normal HTTP request (not WebSocket upgrade) - should fail silently
	req := httptest.NewRequest(http.MethodGet, "/ws", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// The handler should return a non-200 because it can't upgrade
	if w.Code == http.StatusOK {
		t.Error("expected non-200 status for non-WebSocket request")
	}
}
