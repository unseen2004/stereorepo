package google

import (
	"sync"
	"golang.org/x/oauth2"
)

type TokenStore struct {
	tokens map[string]*oauth2.Token
	mu     sync.RWMutex
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]*oauth2.Token),
	}
}

func (s *TokenStore) SetToken(key string, token *oauth2.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[key] = token
}

func (s *TokenStore) GetToken(key string) (*oauth2.Token, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token, ok := s.tokens[key]
	return token, ok
}
