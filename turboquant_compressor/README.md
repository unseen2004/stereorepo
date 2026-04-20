# turboquant_compressor

Efficient LLM quantization tool designed to compress high-dimensional vectors (KV caches and embeddings) down to 2-4 bits using the TurboQuant algorithm.

## Features

- Multi-threaded compression (f32, f16, bf16 to 2-4 bits).
- Full decompression back to f32.
- Batch processing for entire model directories.
- Python extension (PyO3) for direct integration.
- Minimal accuracy loss via random rotation and MSE optimization.

## Installation

Ensure you have Rust and Cargo installed.

```bash
cargo build --release
```

## Usage

### Compression

Compress a single file or an entire directory:

```bash
./target/release/turboquant_compressor compress <input_path> --bits 4 --seed 42
```

### Decompression

Decompress a .tq.safetensors file back to .f32.safetensors:

```bash
./target/release/turboquant_compressor decompress <input_path>.tq.safetensors
```

## Testing

Run the end-to-end verification test:

```bash
cargo test
```

## Python Integration

The Python extension is located in the `python_ext` directory. It can be built using `maturin`.

```bash
cd python_ext
maturin develop
```

Example usage in Python:

```python
import turboquant_compressor
import numpy as np

data = np.random.randn(10, 128).astype(np.float32)
quantized_bytes, norms = turboquant_compressor.compress(data.flatten().tolist(), 128, bits=4)
```
