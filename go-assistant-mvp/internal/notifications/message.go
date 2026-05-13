package notifications

type Notification struct {
	Type           string  `json:"type"`
	Title          string  `json:"title"`
	Body           string  `json:"body"`
	TaskID         string  `json:"task_id,omitempty"`
	DistanceMeters float64 `json:"distance_meters,omitempty"`
}
