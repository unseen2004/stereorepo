package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub_BroadcastToUser(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "?user_id=test-user"
	
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer ws.Close()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	expected := Notification{
		Type:  "test",
		Title: "Hello",
		Body:  "World",
	}
	msg, _ := json.Marshal(expected)
	hub.BroadcastToUser("test-user", msg)

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}

	var actual Notification
	if err := json.Unmarshal(message, &actual); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if actual.Type != expected.Type || actual.Title != expected.Title {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}
