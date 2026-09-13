package main

import "testing"

type person struct {
	// I AM NOT DONE
}

func newPerson(name string, age int) person {
	// I AM NOT DONE
	return person{}
}

func TestPerson(t *testing.T) {
	p := newPerson("Ada", 36)
	if p.name != "Ada" {
		t.Fatalf("p.name = %q, want %q", p.name, "Ada")
	}
	if p.age != 36 {
		t.Fatalf("p.age = %d, want 36", p.age)
	}
}
