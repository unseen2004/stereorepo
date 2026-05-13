package ai

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Config struct {
	APIKey string
}

type TaskSummary struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Status              string  `json:"status"`
	DueAt               string  `json:"due_at,omitempty"`
	HasLocationTrigger  bool    `json:"has_location_trigger"`
	DistanceMeters      float64 `json:"distance_meters,omitempty"`
}

type ContextRequest struct {
	UserLat   float64       `json:"lat"`
	UserLng   float64       `json:"lng"`
	Tasks     []TaskSummary `json:"tasks"`
	TimeOfDay string        `json:"time_of_day"`
}

type Service struct {
	client *genai.Client
	config Config
}

func NewService(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.APIKey == "" {
		return &Service{client: nil, config: cfg}, nil
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, err
	}
	return &Service{client: client, config: cfg}, nil
}

func (s *Service) IsConfigured() bool {
	return s.client != nil
}

func (s *Service) GetSuggestions(ctx context.Context, req ContextRequest) (string, error) {
	if !s.IsConfigured() {
		return "[DEMO MODE - Gemini API key not set] You have " + strconv.Itoa(len(req.Tasks)) + " tasks to complete. The nearest task with a location trigger should be your priority.", nil
	}

	// In 2026, gemini-2.0-flash is deprecated. Using gemini-2.5-flash instead to fulfill the request.
	model := s.client.GenerativeModel("gemini-2.5-flash")
	
	systemPrompt := "You are a task assistant. ALWAYS respond in English. Be specific and concise. Maximum 3 suggestions."
	
	userPrompt := "User has the following pending tasks:\n"
	for _, t := range req.Tasks {
		dueInfo := ""
		if t.DueAt != "" {
			dueInfo = fmt.Sprintf(" (due: %s)", t.DueAt)
		}
		userPrompt += fmt.Sprintf("- %s%s\n", t.Title, dueInfo)
	}
	userPrompt += fmt.Sprintf("\nUser location: lat=%f, lng=%f\nCurrent time: %s\nSuggest what to do next and why.", 
		req.UserLat, req.UserLng, req.TimeOfDay)

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "No suggestions available.", nil
	}

	part := resp.Candidates[0].Content.Parts[0]
	if text, ok := part.(genai.Text); ok {
		return string(text), nil
	}

	return "Could not parse AI response.", nil
}
