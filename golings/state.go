package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type state struct {
	Done []string `json:"done"`
}

func statePath(root string) string {
	return filepath.Join(root, ".golings", "state.json")
}

func loadState(root string) *state {
	data, err := os.ReadFile(statePath(root))
	if err != nil {
		return &state{}
	}
	var s state
	if json.Unmarshal(data, &s) != nil {
		return &state{}
	}
	return &s
}

func (s *state) isDone(ex Exercise) bool {
	for _, name := range s.Done {
		if name == ex.Name {
			return true
		}
	}
	return false
}

func (s *state) markDone(ex Exercise) {
	if !s.isDone(ex) {
		s.Done = append(s.Done, ex.Name)
	}
}

func (s *state) save(root string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	path := statePath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (s *state) doneCount() int {
	count := 0
	for _, ex := range exercises {
		if s.isDone(ex) {
			count++
		}
	}
	return count
}

func (s *state) nextPending() *Exercise {
	for i := range exercises {
		if !s.isDone(exercises[i]) {
			return &exercises[i]
		}
	}
	return nil
}
