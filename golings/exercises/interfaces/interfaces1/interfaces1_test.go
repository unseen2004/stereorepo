package main

import "testing"

type speaker interface {
	speak() string
}

type dog struct{}

type cat struct{}

// I AM NOT DONE

func TestSpeak(t *testing.T) {
	var s speaker = dog{}
	if got := s.speak(); got != "Woof" {
		t.Fatalf("dog.speak() = %q, want %q", got, "Woof")
	}
	s = cat{}
	if got := s.speak(); got != "Meow" {
		t.Fatalf("cat.speak() = %q, want %q", got, "Meow")
	}
}
