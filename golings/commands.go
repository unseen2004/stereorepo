package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	colorGreen = "\x1b[32m"
	colorRed   = "\x1b[31m"
	colorBold  = "\x1b[1m"
	colorDim   = "\x1b[2m"
	colorReset = "\x1b[0m"
)

func run(name string) {
	root := mustRoot()
	st := loadState(root)

	ex := st.nextPending()
	if name != "" {
		var found bool
		ex, found = findExercise(name)
		if !found {
			fmt.Fprintf(os.Stderr, "golings: no exercise named %q\n\n", name)
			similar(name)
			os.Exit(1)
		}
	}
	if ex == nil {
		printStatus(st)
		return
	}

	fmt.Printf("Running %s (%s)\n\n", ex.Name, ex.Mode)
	ok, output := verify(root, *ex)
	printResult(*ex, ok, output)

	if ok {
		st.markDone(*ex)
		if err := st.save(root); err != nil {
			fmt.Printf("warning: could not save state: %v\n", err)
		}
		printStatus(st)
		os.Exit(0)
	}
	os.Exit(1)
}

func hint(name string) {
	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: golings hint <exercise-name>")
		os.Exit(1)
	}
	ex, found := findExercise(name)
	if !found {
		fmt.Fprintf(os.Stderr, "golings: no exercise named %q\n\n", name)
		similar(name)
		os.Exit(1)
	}
	fmt.Printf("%sHint for %s:%s\n\n%s\n", colorBold, ex.Name, colorReset, ex.Hint)
}

func list() {
	root := mustRoot()
	st := loadState(root)
	for i, ex := range exercises {
		mark := "✗"
		if st.isDone(ex) {
			mark = "✓"
		}
		color := colorRed
		if st.isDone(ex) {
			color = colorGreen
		}
		fmt.Printf("%3d. %-14s [%-4s] %s%s%s\n", i+1, ex.Name, ex.Mode, color, mark, colorReset)
	}
	fmt.Println()
	printStatus(st)
}

func verifyAll() {
	root := mustRoot()
	var failed []string
	for _, ex := range exercises {
		ok, output := verify(root, ex)
		if ok {
			fmt.Printf("%s✓ %s%s\n", colorGreen, ex.Name, colorReset)
		} else {
			fmt.Printf("%s✗ %s%s\n", colorRed, ex.Name, colorReset)
			failed = append(failed, ex.Name)
			if output != "" {
				for _, line := range strings.Split(output, "\n") {
					fmt.Printf("    %s\n", line)
				}
			}
		}
	}
	fmt.Printf("\n%d/%d exercises pass\n", len(exercises)-len(failed), len(exercises))
	if len(failed) > 0 {
		os.Exit(1)
	}
}

func reset() {
	root := mustRoot()
	if err := os.RemoveAll(filepathJoin(root, ".golings")); err != nil {
		fmt.Fprintf(os.Stderr, "golings: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Progress reset.")
}

func similar(name string) {
	var close []string
	for _, ex := range exercises {
		if strings.Contains(ex.Name, name) || strings.Contains(name, ex.Name) {
			close = append(close, ex.Name)
		}
	}
	sort.Strings(close)
	if len(close) > 0 {
		fmt.Println("Did you mean one of these?")
		for _, n := range close {
			fmt.Println("  " + n)
		}
	}
}

func filepathJoin(elem ...string) string {
	return strings.Join(elem, string(os.PathSeparator))
}
