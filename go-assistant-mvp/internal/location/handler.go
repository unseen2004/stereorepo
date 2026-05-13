package location

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RecordLocation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string  `json:"user_id"`
		Lat    float64 `json:"lat"`
		Lng    float64 `json:"lng"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RecordLocation(body.UserID, Point{Lat: body.Lat, Lng: body.Lng}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "recorded"})
}

func (h *Handler) GetNearby(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Lat    float64 `json:"lat"`
		Lng    float64 `json:"lng"`
		Radius float64 `json:"radius"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Radius == 0 {
		body.Radius = 300 // default 300m
	}

	tasks, err := h.service.GetNearbyTasks(Point{Lat: body.Lat, Lng: body.Lng}, body.Radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) CheckProximity(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Current Point   `json:"current"`
		Target  Point   `json:"target"`
		Radius  float64 `json:"radius"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	near := h.service.IsNearPoint(body.Current, body.Target, body.Radius)
	dist := h.service.GetDistance(body.Current, body.Target)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"near":            near,
		"distance_meters": dist,
	})
}
