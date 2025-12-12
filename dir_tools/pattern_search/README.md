# Pattern Search (rgrep)

Fast, idiomatic Rust grep-like CLI with streaming IO and rich feature flags.

---

## Overview

A feature-rich grep-like tool built with Rust, offering advanced search capabilities with multiple output formats and filtering options. Designed for performance with streaming I/O.

---

## Building

```bash
cargo build --release
```

---

## Running

Basic search:
```bash
cargo run -- <pattern> [files...]
```

With options:
```bash
cargo run -- <pattern> [files...] [OPTIONS]
```

---

## Features

- Fast pattern matching with streaming I/O
- Multiple output formats (JSON, plain text)
- Case-sensitive/insensitive search
- Recursive directory search
- Color-coded output
- Verbosity control
- Human-friendly error messages
- Comprehensive test suite

---

## Options

Use `--help` to see all available options:
```bash
cargo run -- --help
```

---

## Testing

Run tests:
```bash
cargo test
```

---

## Author

unseen2004

---

## License

MIT
