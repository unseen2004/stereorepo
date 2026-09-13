package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func watch() {
	root := mustRoot()
	st := loadState(root)

	if printStatus(st) {
		return
	}

	ex := st.nextPending()
	fmt.Printf("\nWatching %s%s%s\n\n", colorDim, exerciseDir(ex), colorReset)
	fmt.Printf("    %sCurrent exercise: %s%s\n", colorBold, ex.Name, colorReset)
	fmt.Println("    Edit the file in your editor and save to check your work.")
	fmt.Println("    Press Ctrl+C to quit.")

	mtimes := snapshot(root, *ex)
	for {
		time.Sleep(400 * time.Millisecond)
		current := snapshot(root, *ex)
		if mtimesEqual(mtimes, current) {
			continue
		}
		mtimes = current

		ok, output := verify(root, *ex)
		if !ok {
			printResult(*ex, false, output)
			continue
		}

		st.markDone(*ex)
		if err := st.save(root); err != nil {
			fmt.Printf("warning: could not save state: %v\n", err)
		}
		printResult(*ex, true, output)

		if printStatus(st) {
			return
		}

		ex = st.nextPending()
		mtimes = snapshot(root, *ex)
		fmt.Printf("\nWatching %s%s%s\n\n", colorDim, exerciseDir(ex), colorReset)
		fmt.Printf("    %sCurrent exercise: %s%s\n", colorBold, ex.Name, colorReset)
	}
}

func printStatus(st *state) bool {
	done := st.doneCount()
	total := len(exercises)
	fmt.Printf("\nProgress: %s%d/%d%s (%d%%)\n", colorBold, done, total, colorReset, done*100/total)
	if done == total {
		fmt.Printf("\n%s🎉 All exercises completed! You're done. 🎉%s\n", colorGreen, colorReset)
		return true
	}
	return false
}

func exerciseDir(ex *Exercise) string {
	return filepath.Join("exercises", filepath.FromSlash(ex.Dir))
}

func snapshot(root string, ex Exercise) map[string]time.Time {
	times := map[string]time.Time{}
	dir := filepath.Join(root, "exercises", filepath.FromSlash(ex.Dir))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return times
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".go" {
			continue
		}
		info, err := e.Info()
		if err == nil {
			times[e.Name()] = info.ModTime()
		}
	}
	return times
}

func mtimesEqual(a, b map[string]time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	for name, t := range a {
		if !b[name].Equal(t) {
			return false
		}
	}
	return true
}
