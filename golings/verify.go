package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func verify(root string, ex Exercise) (bool, string) {
	var cmd *exec.Cmd
	switch ex.Mode {
	case "test":
		cmd = exec.Command("go", "test", "./exercises/"+ex.Dir)
	default:
		cmd = exec.Command("go", "run", "./exercises/"+ex.Dir)
	}
	cmd.Dir = root
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return false, cleanOutput(out.String())
	}

	dir := filepath.Join(root, "exercises", filepath.FromSlash(ex.Dir))
	if hasMarker(dir) {
		return false, "the code compiles and runs, but the `// I AM NOT DONE` marker is still present.\nRemove that comment from the exercise file and you're done."
	}
	return true, strings.TrimSpace(out.String())
}

// "I AM NOT DONE" comment.
func hasMarker(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err == nil && bytes.Contains(data, []byte(notDoneMarker)) {
			return true
		}
	}
	return false
}

func cleanOutput(out string) string {
	out = strings.TrimSpace(out)
	out = strings.TrimSuffix(out, "exit status 1")
	out = strings.TrimSpace(out)
	return out
}

func printResult(ex Exercise, ok bool, output string) {
	if ok {
		fmt.Printf("%s✓ exercise %s passed%s\n", colorGreen, ex.Name, colorReset)
		if output != "" {
			fmt.Println(output)
		}
	} else {
		fmt.Printf("%s✗ exercise %s failed%s\n", colorRed, ex.Name, colorReset)
		if output != "" {
			fmt.Println(output)
		}
		fmt.Printf("Need a hint? Run: %sgolings hint %s%s\n", colorDim, ex.Name, colorReset)
	}
}
