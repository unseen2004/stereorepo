package main

import (
	"testing"
	"time"
)

func send(c chan int, n int) {
	// I AM NOT DONE
}

func TestSend(t *testing.T) {
	c := make(chan int)
	go send(c, 42)
	select {
	case got := <-c:
		if got != 42 {
			t.Fatalf("received %d, want 42", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out - nothing was sent on the channel")
	}
}
