# Git Visualizer

Terminal-based contribution graph for local git repositories.

## Installation

```bash
go build -o git-visualizer
```

## Usage

### Add repositories
Scan a directory for git repositories to track:
```bash
./git-visualizer -add /path/to/projects
```

### View stats
Show contribution graph for a specific email:
```bash
./git-visualizer -email user@example.com
```
