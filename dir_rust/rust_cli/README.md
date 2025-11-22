# Rust CLI

Command-line interface utilities and tools in Rust.

---

## Overview

Collection of CLI utilities demonstrating Rust command-line application development with proper argument parsing, JSON serialization, and error handling.

---

## Building

```bash
cargo build --release
```

---

## Running

```bash
cargo run -- [COMMAND] [OPTIONS]
```

---

## Features

- Command-line argument parsing with clap
- JSON serialization/deserialization with serde
- Error handling with anyhow
- Collection utilities and macros
- Structured CLI with subcommands

---

## Dependencies

- anyhow: Error handling
- clap: CLI parsing with derive macros
- serde: Serialization framework
- serde_json: JSON support
- collection_macros: Collection utilities

---

## Usage

Use `--help` to see available commands and options:
```bash
cargo run -- --help
```
