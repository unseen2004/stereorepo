package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateRoundTrip(t *testing.T) {
	root := t.TempDir()
	s := &state{}
	if s.nextPending() == nil || s.nextPending().Name != "hello1" {
		t.Fatal("fresh state should point at hello1")
	}
	s.markDone(*s.nextPending())
	if got := s.doneCount(); got != 1 {
		t.Fatalf("doneCount = %d, want 1", got)
	}
	if s.nextPending().Name != "variables1" {
		t.Fatalf("next pending = %s, want variables1", s.nextPending().Name)
	}
	if err := s.save(root); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(statePath(root)); err != nil {
		t.Fatalf("state file not written: %v", err)
	}
	reloaded := loadState(root)
	if reloaded.doneCount() != 1 || !reloaded.isDone(exercises[0]) {
		t.Fatalf("reloaded state wrong: %+v", reloaded.Done)
	}
}

func TestStateAllDone(t *testing.T) {
	s := &state{}
	for i := range exercises {
		s.markDone(exercises[i])
	}
	if s.nextPending() != nil {
		t.Fatal("nextPending should be nil when everything is done")
	}
	if got := s.doneCount(); got != len(exercises) {
		t.Fatalf("doneCount = %d, want %d", got, len(exercises))
	}
}

func TestHasMarker(t *testing.T) {
	root := t.TempDir()
	with := filepath.Join(root, "with")
	without := filepath.Join(root, "without")
	os.MkdirAll(with, 0o755)
	os.MkdirAll(without, 0o755)
	os.WriteFile(filepath.Join(with, "main.go"), []byte("package main\n// I AM NOT DONE\nfunc main() {}\n"), 0o644)
	os.WriteFile(filepath.Join(without, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644)

	if !hasMarker(with) {
		t.Fatal("hasMarker should find the marker")
	}
	if hasMarker(without) {
		t.Fatal("hasMarker should not find a marker that is not there")
	}
	if hasMarker(filepath.Join(root, "missing")) {
		t.Fatal("hasMarker should return false for missing directories")
	}
}

func TestFindExercise(t *testing.T) {
	if _, ok := findExercise("channels2"); !ok {
		t.Fatal("expected channels2 to exist")
	}
	if _, ok := findExercise("doesnotexist"); ok {
		t.Fatal("expected doesnotexist to not exist")
	}
}

func TestExercisesWellFormed(t *testing.T) {
	names := map[string]bool{}
	for _, ex := range exercises {
		if names[ex.Name] {
			t.Fatalf("duplicate exercise name %q", ex.Name)
		}
		names[ex.Name] = true
		if ex.Mode != "run" && ex.Mode != "test" {
			t.Fatalf("%s: unknown mode %q", ex.Name, ex.Mode)
		}
		if ex.Hint == "" {
			t.Fatalf("%s: missing hint", ex.Name)
		}
		if ex.Dir == "" {
			t.Fatalf("%s: missing dir", ex.Name)
		}
	}
	if len(exercises) == 0 {
		t.Fatal("no exercises defined")
	}
}
