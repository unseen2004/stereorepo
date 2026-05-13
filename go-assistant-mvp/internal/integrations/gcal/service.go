package gcal

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Service struct {
	config   Config
	OAuthCfg *oauth2.Config
}

func NewService(cfg Config) *Service {
	oc := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{calendar.CalendarReadonlyScope, calendar.CalendarEventsScope, gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
	}
	return &Service{config: cfg, OAuthCfg: oc}
}

func (s *Service) GetAuthURL(state string) string {
	return s.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.OAuthCfg.Exchange(ctx, code)
}

func (s *Service) GetTodayEvents(ctx context.Context, token *oauth2.Token) ([]*calendar.Event, error) {
	client := s.OAuthCfg.Client(ctx, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tMin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	tMax := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Format(time.RFC3339)

	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(tMin).
		TimeMax(tMax).
		MaxResults(20).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

func (s *Service) EventsToTasks(events []*calendar.Event) []map[string]interface{} {
	tasks := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		t := map[string]interface{}{
			"title":       e.Summary,
			"description": e.Description,
			"status":      "pending",
			"source":      "google_calendar",
			"external_id": e.Id,
		}
		if e.Start.DateTime != "" {
			t["due_at"] = e.Start.DateTime
		} else {
			t["due_at"] = e.Start.Date
		}
		tasks = append(tasks, t)
	}
	return tasks
}
