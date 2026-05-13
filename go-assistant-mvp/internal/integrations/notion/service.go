package notion

import (
	"context"

	"github.com/jomei/notionapi"
)

type Config struct {
	APIKey     string
	DatabaseID string
}

type Service struct {
	client *notionapi.Client
	config Config
}

func NewService(cfg Config) *Service {
	client := notionapi.NewClient(notionapi.Token(cfg.APIKey))
	return &Service{client: client, config: cfg}
}

func (s *Service) QueryDatabase(ctx context.Context) ([]notionapi.Page, error) {
	resp, err := s.client.Database.Query(ctx, notionapi.DatabaseID(s.config.DatabaseID), nil)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (s *Service) PageToTask(page notionapi.Page) map[string]interface{} {
	title := ""

	// Look for any property of type "title"
	for _, prop := range page.Properties {
		if prop.GetType() == notionapi.PropertyTypeTitle {
			if titleProp, ok := prop.(*notionapi.TitleProperty); ok && len(titleProp.Title) > 0 {
				title = titleProp.Title[0].PlainText
			}
			break
		}
	}

	// Fallback to page ID if title is empty
	if title == "" {
		title = page.ID.String()
	}

	return map[string]interface{}{
		"title":       title,
		"status":      "pending",
		"source":      "notion",
		"external_id": page.ID.String(),
		"description": page.URL,
		"due_at":      nil,
	}
}
