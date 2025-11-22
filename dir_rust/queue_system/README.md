# Queue System

QR Code Queue Management System in Rust.

---

## Overview

Queue management application with QR code generation for tickets and OpenCV-based scanning. Features automatic queue progression and real-time status checking.

---

## Building

```bash
cargo build --release
```

---

## Running

Interactive mode:
```bash
cargo run -- run
```

Generate QR code:
```bash
cargo run -- generate
```

Scan QR code with camera:
```bash
cargo run -- scan
```

Check queue status:
```bash
cargo run -- status
```

---

## Features

- Generate unique QR code tickets
- Camera-based QR code scanning with OpenCV
- Auto-remove first person after 20 seconds
- Re-scan shows current position in queue
- Terminal-based CLI interface
- Real-time queue status

---

## Dependencies

- qrcode: QR code generation
- opencv: Camera and QR scanning
- tokio: Async runtime
- clap: CLI argument parsing
- crossterm: Terminal UI
- chrono: Time management
- uuid: Unique identifiers

