# FISML - File Integrity and Secret Leakage Monitor

Minimal Rust MVP for monitoring file changes and scanning for potential secret leaks.

---

## Features

- Baseline scan of configured paths
- Detect new/modified files using hash comparison
- Optional secret scanning with regex and entropy heuristic
- SQLite-backed state and append-only event chain
- Watch mode with graceful shutdown
- Event hash chain integrity verification

---

## Installation

Build from source:
```bash
cargo build --release
```

---

## Usage

Initialize configuration:
```bash
cargo run -- init-config
```

Run baseline scan with secret detection:
```bash
cargo run -- baseline --secrets
```

List recent events:
```bash
cargo run -- list-events --limit 20
```

Watch for changes continuously:
```bash
cargo run -- watch --secrets
```

Verify event chain integrity:
```bash
cargo run -- verify-chain
```

---

## Configuration

Configuration file (`fisml.toml`) is automatically generated on first run. Edit to configure:
- Monitored paths
- Ignore patterns
- Entropy threshold for secret detection

---

## Logging

Enable verbose logging:
```bash
RUST_LOG=info cargo run -- watch
```

---

## Disclaimer

MVP quality. Not production hardened. Use at your own risk and only on authorized systems.

