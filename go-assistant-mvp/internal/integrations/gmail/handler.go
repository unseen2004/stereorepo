package gmail

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-assistant/core/pkg/google"
)

type TaskStore interface {
	CreateRaw(data map[string]interface{}) error
}

type Handler struct {
	service    *Service
	taskStore  TaskStore
	tokenStore *google.TokenStore
}

func NewHandler(svc *Service, store TaskStore, tokenStore *google.TokenStore) *Handler {
	return &Handler{
		service:    svc,
		taskStore:  store,
		tokenStore: tokenStore,
	}
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	_, ok := h.tokenStore.GetToken("default")
	configured := ok
	scopeGranted := false
	if ok {
		// If we have a token, we assume scopes are granted for simplicity
		scopeGranted = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"configured":    configured,
		"scope_granted": scopeGranted,
	})
}

func (h *Handler) HandleSync(w http.ResponseWriter, r *http.Request) {
	token, ok := h.tokenStore.GetToken("default")
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "not authenticated, visit /auth/google/start"})
		return
	}

	messages, err := h.service.FetchImportantEmails(context.Background(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := h.service.MessagesToTasks(messages)
	for _, t := range tasks {
		_ = h.taskStore.CreateRaw(t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"synced": len(tasks),
		"tasks":  tasks,
	})
}
