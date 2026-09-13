# golings

Interactive exercises for learning Go, in the style of rustlings. Each
exercise is a small Go program with an intentional bug: some don't compile,
some fail their tests, some print the wrong thing. Fix them one by one.

## Setup

Requires Go 1.22+.

```sh
go build -o golings .
./golings
```

The watch loop shows the current exercise and re-checks it every time you
save the file. When it passes, progress is recorded and the next exercise
loads.

Every exercise file starts with a `// I AM NOT DONE` line. Remove it to
finish the exercise.

## Commands

```sh
golings              # start watch mode
golings watch        # same as above
golings run [name]   # run the next pending exercise, or a named one
golings hint <name>  # show a hint for an exercise
golings list         # list all exercises with their status
golings verify       # run every exercise and report results
golings reset        # clear all progress
```

Run the tool from anywhere inside the repo; it finds the repo root via
`go.mod`. Progress is stored in `.golings/state.json`.

## The course

30 exercises in `exercises/<section>/<name>/`:

- `hello1` — first program
- `variables1..5` — declarations, `:=`, constants, shadowing, zero values
- `types1..3` — numeric and string conversions
- `functions1..4` — returns, variadics, closures
- `controlflow1..3` — if/else, for, range
- `strings1..2` — runes, iterating strings
- `structs1..2` — struct literals, methods
- `slices1..2` — append, make/copy
- `maps1`, `pointers1`, `interfaces1`
- `errors1..2` — returning and formatting errors
- `generics1` — type parameters
- `channels1..2` — goroutines and channels

Exercise modes:

- `run` — the file must compile and exit cleanly (`go run`)
- `test` — the tests in the file must pass (`go test`)

## Adding exercises

1. Create the file in `exercises/<section>/<name>/`:
   - run mode: `main.go` with `package main` and `func main()`
   - test mode: `<name>_test.go` with `package main` and your tests
2. Add the `// I AM NOT DONE` marker.
3. Register it in `exercise.go` with a name, dir, mode, and hint.

Exercises are regular Go files, no DSL involved.
