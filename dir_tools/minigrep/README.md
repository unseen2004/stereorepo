# Minigrep

Simple grep-like text search tool implementation in Rust.

---

## Overview

A basic implementation of the grep command-line utility for searching text patterns in files. Created as a learning project to practice Rust fundamentals.

---

## Building

```bash
cargo build --release
```

---

## Running

Search for a pattern in a file:
```bash
cargo run -- <pattern> <file>
```

Or use the compiled binary:
```bash
./target/release/minigrep <pattern> <file>
```

---

## Features

- Text pattern searching
- File reading and processing
- Command-line argument parsing
- Case-sensitive/insensitive search
- Error handling

---

## Example

```bash
cargo run -- "search_term" example.txt
```
