package ai

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-assistant/core/internal/tasks"
)

type TaskStore interface {
	GetAll() ([]tasks.Task, error)
}

type Handler struct {
	service   *Service
	taskStore TaskStore
}

func NewHandler(svc *Service, store TaskStore) *Handler {
	return &Handler{service: svc, taskStore: store}
}

func (h *Handler) HandleSuggest(w http.ResponseWriter, r *http.Request) {
	var req ContextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch tasks from TaskStore
	dbTasks, err := h.taskStore.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	// Map DB tasks to TaskSummary for the service
	req.Tasks = make([]TaskSummary, 0, len(dbTasks))
	for _, t := range dbTasks {
		if t.Status == "pending" {
			summary := TaskSummary{
				ID:     t.ID,
				Title:  t.Title,
				Status: t.Status,
			}
			if t.DueAt != nil {
				summary.DueAt = t.DueAt.Format(time.RFC3339)
			}
			if t.LocationTrigger != nil {
				summary.HasLocationTrigger = true
			}
			req.Tasks = append(req.Tasks, summary)
		}
	}

	if req.TimeOfDay == "" {
		req.TimeOfDay = time.Now().UTC().Format(time.RFC3339)
	}

	suggestions, err := h.service.GetSuggestions(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions":   suggestions,
		"ai_configured": h.service.IsConfigured(),
		"task_count":    len(req.Tasks),
	})
}

func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"configured": h.service.IsConfigured(),
		"model":      "gemini-2.0-flash", // Returning as requested by user
	})
}
