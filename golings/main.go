package main

import (
	"fmt"
	"os"
)

const helpText = `golings - interactive exercises for learning Go (rustlings-style)

Usage:
  golings               start watch mode (re-checks the current exercise on save)
  golings watch         same as running with no arguments
  golings run [name]    run the next pending exercise, or the named one
  golings hint <name>   show a hint for an exercise
  golings list          list all exercises and their status
  golings verify        run all exercises and report which ones pass
  golings reset         clear all progress

Exercises live in ./exercises/<section>/<name>/. Edit the file in place,
save it, and the watch loop verifies your solution automatically.
`

func main() {
	cmd := "watch"
	name := ""
	args := os.Args[1:]
	if len(args) > 0 {
		cmd = args[0]
		if len(args) > 1 {
			name = args[1]
		}
	}
	switch cmd {
	case "watch":
		watch()
	case "run":
		run(name)
	case "hint":
		hint(name)
	case "list":
		list()
	case "verify":
		verifyAll()
	case "reset":
		reset()
	case "help", "-h", "--help":
		fmt.Print(helpText)
	default:
		fmt.Fprintf(os.Stderr, "golings: unknown command %q\n\n%s", cmd, helpText)
		os.Exit(1)
	}
}
