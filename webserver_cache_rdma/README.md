# C++ Multithreaded Web Server with LRU Cache and RDMA

High-performance HTTP/1.1 static file server with in-memory caching and optional RDMA fast path.

---

## Overview

Production-quality static web server written in modern C++ featuring asynchronous networking, intelligent caching, and optional RDMA support for ultra-low latency data delivery.

---

## Features

**HTTP/1.1:**
- GET and HEAD methods
- Keep-alive and request pipelining
- MIME type detection
- Path traversal protection

**Caching:**
- Thread-safe in-memory LRU cache
- Configurable size limits
- ETag and Last-Modified support

**RDMA (optional):**
- rdma_cm + ibverbs integration
- Custom binary protocol over SEND/RECV
- Shared cache with HTTP path
- Pre-posted receives per connection

**Operational:**
- Clean shutdown on signals
- Metrics endpoint (/metrics)
- Docker packaging

---

## Building

Prerequisites: CMake 3.16+, C++17 compiler, Boost.System

Build:
```bash
mkdir build
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j
```

With RDMA:
```bash
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release -DENABLE_RDMA=ON
cmake --build build -j
```

---

## Running

Basic HTTP server:
```bash
./build/webserver --port 8080 --threads 4 --doc-root ./public --cache.mem-mb 128
```

With RDMA:
```bash
./build/webserver --port 8080 --doc-root ./public --rdma.enable --rdma.port 7471
```

---

## Docker

Build image:
```bash
docker build -t webserver:latest .
```

Run container:
```bash
docker run --rm -p 8080:8080 -v $(pwd)/public:/site webserver:latest \
  /app/webserver --port 8080 --doc-root /site --cache.mem-mb 128
```

---

## Configuration

**HTTP Options:**
- `--port N` - HTTP port (default 8080)
- `--threads N` - Worker threads (0 = auto)
- `--doc-root PATH` - Document root (default ./public)
- `--cache.mem-mb N` - Cache size in MB (default 128)
- `--keepalive-timeout-ms N` - Keep-alive timeout (default 10000)

**RDMA Options:**
- `--rdma.enable` - Enable RDMA endpoint
- `--rdma.bind IP` - Bind address (default 0.0.0.0)
- `--rdma.port N` - RDMA port (default 7471)
- `--rdma.recv-bufs N` - Receive buffers (default 64)
- `--rdma.send-chunk N` - Send chunk size (default 32768)

---

## Project Structure

```
webserver_cache_rdma/
├── src/
│   ├── main.cpp              # Entry point
│   ├── server.{hpp,cpp}      # HTTP server
│   ├── session.{hpp,cpp}     # HTTP session
│   ├── http/                 # HTTP parsing and response
│   ├── cache/                # LRU cache implementation
│   ├── fs/                   # File system utilities
│   ├── rdma/                 # RDMA implementation
│   └── util/                 # Configuration, logging, metrics
├── CMakeLists.txt            # Build configuration
├── Dockerfile                # Container image
└── public/                   # Default document root
```

---

## Metrics

Access metrics at `/metrics`:
```bash
curl http://localhost:8080/metrics
```

Includes:
- Request counters
- Response status counts
- Cache hit/miss statistics
- Bytes served
- RDMA operation counts (if enabled)

---

## RDMA Protocol

Binary protocol over SEND/RECV:

Request:
- Header: `{uint8 op, uint16 path_len}`
- Op=1 (GET): followed by path string
- Op=2 (PING): no payload

Response:
- Header: `{uint16 status, uint64 content_len, uint32 chunk_size}`
- Followed by content in chunks

---

## Performance Tips

- Increase `--threads` for multi-core systems
- Size `--cache.mem-mb` to hold frequently accessed files
- Use RDMA for trusted internal networks requiring lowest latency
- Tune `--rdma.recv-bufs` and `--rdma.send-chunk` for workload

---

## Security

- Static file serving only (no directory listings)
- Path traversal protection
- RDMA endpoint for trusted networks only
- Consider network ACLs or proxy-level mTLS

---

## Dependencies

- Boost.Asio (networking)
- fmt (formatting)
- rdma-core (optional, for RDMA support)

---

## License

Choose appropriate license for your project (e.g., MIT, Apache-2.0)
