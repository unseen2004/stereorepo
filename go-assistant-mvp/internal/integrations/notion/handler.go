package notion

import (
	"context"
	"encoding/json"
	"net/http"
)

type TaskStore interface {
	CreateRaw(data map[string]interface{}) error
}

type Handler struct {
	service   *Service
	taskStore TaskStore
}

func NewHandler(svc *Service, store TaskStore) *Handler {
	return &Handler{service: svc, taskStore: store}
}

func (h *Handler) HandleSync(w http.ResponseWriter, r *http.Request) {
	if h.service.config.APIKey == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "Notion API key not configured. Set NOTION_API_KEY env var"})
		return
	}
	if h.service.config.DatabaseID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "Notion database ID not configured. Set NOTION_DATABASE_ID env var"})
		return
	}

	pages, err := h.service.QueryDatabase(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := make([]map[string]interface{}, 0, len(pages))
	for _, page := range pages {
		taskData := h.service.PageToTask(page)
		tasks = append(tasks, taskData)
		if h.taskStore != nil {
			_ = h.taskStore.CreateRaw(taskData)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"synced": len(tasks),
		"tasks":  tasks,
	})
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	configured := h.service.config.APIKey != "" && h.service.config.DatabaseID != ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"configured":      configured,
		"api_key_set":     h.service.config.APIKey != "",
		"database_id_set": h.service.config.DatabaseID != "",
	})
}
