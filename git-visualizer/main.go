package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	outOfRange      = 99999
	daysInPeriod    = 183
	weeksInPeriod   = 26
)

type column []int

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "Folder to scan for git repositories")
	flag.StringVar(&email, "email", "", "Email to filter commits")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	if email == "" {
		fmt.Println("Please provide an email using -email")
		return
	}

	stats(email)
}

func scan(path string) {
	fmt.Printf("Scanning %s for repositories...\n", path)
	repos := findRepos(path)
	storeRepos(repos)
	fmt.Printf("Added %d new repositories.\n", len(repos))
}

func findRepos(root string) []string {
	var repos []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && info.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
			return filepath.SkipDir
		}
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "vendor") {
			return filepath.SkipDir
		}
		return nil
	})
	return repos
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".git-visualizer")
}

func storeRepos(newRepos []string) {
	path := getConfigPath()
	existing := loadRepos(path)
	
	unique := make(map[string]bool)
	for _, r := range existing {
		unique[r] = true
	}
	for _, r := range newRepos {
		unique[r] = true
	}

	var all []string
	for r := range unique {
		all = append(all, r)
	}

	os.WriteFile(path, []byte(strings.Join(all, "\n")), 0644)
}

func loadRepos(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var repos []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			repos = append(repos, line)
		}
	}
	return repos
}

func stats(email string) {
	commits := collectCommits(email)
	displayStats(commits)
}

func collectCommits(email string) map[int]int {
	repos := loadRepos(getConfigPath())
	commits := make(map[int]int)
	for i := 0; i <= daysInPeriod; i++ {
		commits[i] = 0
	}

	for _, path := range repos {
		fillCommits(email, path, commits)
	}
	return commits
}

func fillCommits(email, path string, commits map[int]int) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return
	}

	ref, err := repo.Head()
	if err != nil {
		return
	}

	it, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return
	}

	offset := getWeekdayOffset()
	it.ForEach(func(c *object.Commit) error {
		daysAgo := calculateDaysAgo(c.Author.When) + offset
		if c.Author.Email == email && daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
}

func calculateDaysAgo(t time.Time) int {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	target := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	
	days := int(today.Sub(target).Hours() / 24)
	if days < 0 || days > daysInPeriod {
		return outOfRange
	}
	return days
}

func getWeekdayOffset() int {
	return int(time.Now().Weekday())
}

func displayStats(commits map[int]int) {
	keys := make([]int, 0, len(commits))
	for k := range commits {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	cols := make(map[int]column)
	for _, k := range keys {
		week := k / 7
		day := k % 7
		if _, ok := cols[week]; !ok {
			cols[week] = make(column, 7)
		}
		cols[week][day] = commits[k]
	}

	printMonths()
	for day := 0; day < 7; day++ {
		printDayLabel(day)
		for week := weeksInPeriod; week >= 0; week-- {
			val := 0
			if col, ok := cols[week]; ok {
				val = col[day]
			}
			isToday := week == 0 && day == getWeekdayOffset()
			renderCell(val, isToday)
		}
		fmt.Println()
	}
}

func renderCell(val int, isToday bool) {
	color := "\033[0;37;30m"
	switch {
	case isToday:
		color = "\033[1;37;45m"
	case val > 10:
		color = "\033[1;30;42m"
	case val >= 5:
		color = "\033[1;30;43m"
	case val > 0:
		color = "\033[1;30;47m"
	}

	format := "  - "
	if val > 0 {
		if val < 10 {
			format = fmt.Sprintf("  %d ", val)
		} else if val < 100 {
			format = fmt.Sprintf(" %d ", val)
		} else {
			format = fmt.Sprintf("%d ", val)
		}
	}
	fmt.Printf("%s%s\033[0m", color, format)
}

func printMonths() {
	now := time.Now()
	start := now.AddDate(0, 0, -daysInPeriod)
	fmt.Print("      ")
	
	currentMonth := start.Month()
	for i := 0; i <= weeksInPeriod; i++ {
		date := start.AddDate(0, 0, i*7)
		if date.Month() != currentMonth || i == 0 {
			fmt.Printf("%s ", date.Month().String()[:3])
			currentMonth = date.Month()
		} else {
			fmt.Print("    ")
		}
	}
	fmt.Println()
}

func printDayLabel(day int) {
	labels := map[int]string{1: "Mon", 3: "Wed", 5: "Fri"}
	if label, ok := labels[day]; ok {
		fmt.Printf(" %s ", label)
	} else {
		fmt.Print("     ")
	}
}
