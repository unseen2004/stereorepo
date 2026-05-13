package location

import ("log"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/go-assistant/core/internal/notifications"
	"github.com/jmoiron/sqlx"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Event struct {
	UserID     string    `db:"user_id"`
	Point      Point     `db:"point"`
	RecordedAt time.Time `db:"recorded_at"`
}

type HubInterface interface {
	BroadcastToUser(string, []byte)
}

type Service struct {
	db  *sqlx.DB
	hub HubInterface
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func NewServiceWithHub(db *sqlx.DB, hub HubInterface) *Service {
	return &Service{db: db, hub: hub}
}

func (s *Service) RecordLocation(userID string, p Point) error {
	query := `INSERT INTO location_events (user_id, lat, lng) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, userID, p.Lat, p.Lng)
	if err != nil {
		return err
	}

	if s.hub == nil {
		return nil
	}

	// Check for nearby tasks
	var tasks []struct {
		ID              string          `db:"id"`
		Title           string          `db:"title"`
		LocationTrigger json.RawMessage `db:"location_trigger"`
	}
	err = s.db.Select(&tasks, "SELECT id, title, location_trigger FROM tasks WHERE location_trigger IS NOT NULL AND status = 'pending'")
	if err != nil {
		return err
	}

	for _, t := range tasks {
		var tp Point
		if err := json.Unmarshal(t.LocationTrigger, &tp); err != nil {
			continue
		}

		if s.IsNearPoint(p, tp, 300) {
			dist := s.GetDistance(p, tp)
			n := notifications.Notification{
				Type:           "task_reminder",
				Title:          "Zadanie w pobliżu!",
				Body:           fmt.Sprintf("%s - jesteś blisko!", t.Title),
				TaskID:         t.ID,
				DistanceMeters: dist,
			}
			nJSON, _ := json.Marshal(n)
			log.Printf("Broadcasting notification for task %s to user %s", t.ID, userID)
			s.hub.BroadcastToUser(userID, nJSON)
		}
	}

	return nil
}

func (s *Service) GetNearbyTasks(p Point, radiusMeters float64) ([]map[string]interface{}, error) {
	var tasks []struct {
		ID              string          `db:"id"`
		Title           string          `db:"title"`
		LocationTrigger json.RawMessage `db:"location_trigger"`
	}
	err := s.db.Select(&tasks, "SELECT id, title, location_trigger FROM tasks WHERE location_trigger IS NOT NULL AND status = 'pending'")
	if err != nil {
		return nil, err
	}

	nearby := make([]map[string]interface{}, 0)
	for _, t := range tasks {
		var tp Point
		if err := json.Unmarshal(t.LocationTrigger, &tp); err != nil {
			continue
		}

		dist := s.GetDistance(p, tp)
		if dist <= radiusMeters {
			nearby = append(nearby, map[string]interface{}{
				"id":              t.ID,
				"title":           t.Title,
				"distance_meters": dist,
			})
		}
	}
	return nearby, nil
}

func (s *Service) IsNearPoint(current Point, target Point, radiusMeters float64) bool {
	return s.GetDistance(current, target) <= radiusMeters
}

func (s *Service) GetDistance(current Point, target Point) float64 {
	const R = 6371000.0 // Earth radius in meters
	
	lat1 := current.Lat * math.Pi / 180
	lng1 := current.Lng * math.Pi / 180
	lat2 := target.Lat * math.Pi / 180
	lng2 := target.Lng * math.Pi / 180

	dLat := lat2 - lat1
	dLng := lng2 - lng1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return R * c
}
