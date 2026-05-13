package tasks

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		h.error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.store.GetByID(id)
	if err != nil {
		h.error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	t.ID = id
	if err := h.store.Update(&t); err != nil {
		h.error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		h.error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Create(&t); err != nil {
		h.error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) error(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
