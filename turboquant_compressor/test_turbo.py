import turbo_compressor
import numpy as np

# 1. Przygotowanie przykładowych danych (np. 10 wektorów o wymiarze 128)
dim = 128
num_vectors = 10
data = np.random.randn(num_vectors, dim).astype(np.float32)

print(f"Oryginalny rozmiar (f32): {data.nbytes} bajtów")

# 2. Kompresja (wymaga spłaszczonej listy/tablicy)
# turbo_compressor.compress(data, dimension, bits)
quantized_bytes, norms = turbo_compressor.compress(data.flatten().tolist(), dim, bits=4)

print(f"Skompresowany rozmiar (indeksy): {len(quantized_bytes)} bajtów")
print(f"Rozmiar norm (f32): {len(norms) * 4} bajtów")
print(f"Całkowity rozmiar: {len(quantized_bytes) + len(norms) * 4} bajtów")
print(f"Stopień kompresji: {data.nbytes / (len(quantized_bytes) + len(norms) * 4):.2f}x")

# 3. Wyświetlenie pierwszych 5 norm
print("\nPierwsze 5 norm wektorów:", norms[:5])
