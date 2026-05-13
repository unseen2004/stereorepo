package gcal

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

func (h *Handler) HandleAuthStart(w http.ResponseWriter, r *http.Request) {
	url := h.service.GetAuthURL("go-assistant")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := h.service.ExchangeCode(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tokenStore.SetToken("default", token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "authenticated"})
}

func (h *Handler) HandleSyncToday(w http.ResponseWriter, r *http.Request) {
	token, ok := h.tokenStore.GetToken("default")
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "not authenticated, visit /auth/google/start"})
		return
	}

	events, err := h.service.GetTodayEvents(context.Background(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := h.service.EventsToTasks(events)
	for _, t := range tasks {
		_ = h.taskStore.CreateRaw(t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
