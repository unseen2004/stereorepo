package gmail

import (
	"context"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Service struct {
	oauthCfg *oauth2.Config
	keywords []string
}

func NewService(oauthCfg *oauth2.Config) *Service {
	return &Service{
		oauthCfg: oauthCfg,
		keywords: []string{
			"faktura", "paczka", "zamówienie", "invoice", "order", "delivery", "reminder", "deadline",
			"PWr", "USOS", "zajecia", "wyniki", "ocena", "platnosc",
			"potwierdzenie", "rejestracja", "payment", "confirmation",
			"meeting", "appointment", "task", "todo",
		},
	}
}

func (s *Service) FetchImportantEmails(ctx context.Context, token *oauth2.Token) ([]*gmail.Message, error) {
	client := s.oauthCfg.Client(ctx, token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Fetch last 20 unread emails
	r, err := srv.Users.Messages.List("me").Q("is:unread").MaxResults(20).Do()
	if err != nil {
		return nil, err
	}

	messages := make([]*gmail.Message, 0)
	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Do()
		if err != nil {
			continue
		}

		if s.isImportant(msg) {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func (s *Service) isImportant(msg *gmail.Message) bool {
	content := strings.ToLower(msg.Snippet)
	subject := ""
	for _, h := range msg.Payload.Headers {
		if h.Name == "Subject" {
			subject = strings.ToLower(h.Value)
			break
		}
	}

	for _, kw := range s.keywords {
		if strings.Contains(content, kw) || strings.Contains(subject, kw) {
			return true
		}
	}
	return false
}

func (s *Service) MessagesToTasks(messages []*gmail.Message) []map[string]interface{} {
	tasks := make([]map[string]interface{}, 0, len(messages))
	for _, m := range messages {
		subject := "No Subject"
		from := "Unknown Sender"
		for _, h := range m.Payload.Headers {
			if h.Name == "Subject" {
				subject = h.Value
			}
			if h.Name == "From" {
				from = h.Value
			}
		}

		t := map[string]interface{}{
			"title":       "Email: " + subject,
			"description": "From: " + from + "\nSnippet: " + m.Snippet,
			"source":      "gmail",
			"external_id": m.Id,
			"status":      "pending",
		}
		tasks = append(tasks, t)
	}
	return tasks
}
