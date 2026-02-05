#pragma once
#include "bitstream.hpp"
#include "tga.hpp"
#include <cstdint>
#include <vector>

// Quantizer class using `double` for precision.
struct Quantizer {
  std::vector<double> reconstruction_values; // Size 2^k

  // Train quantizer on a dataset using the percentile method.
  void train(std::vector<double> &data, int k);

  // Quantize a value to an index (0..2^k-1).
  uint32_t quantize(double val) const;

  // Dequantize index to value.
  double dequantize(uint32_t index) const;
};

struct EncodedData {
  int width, height;
  int k;
  Quantizer quant_L[3]; // R, G, B
  Quantizer quant_H[3]; // R, G, B
  std::vector<uint8_t> bitstream;
};

// Mode 5: Filter (Low/High) + DPCM on Low + Quantize both.
EncodedData encode(const Image &img, int k);
Image decode(const EncodedData &data);

double calculateMSE(const Image &original, const Image &decoded);
double calculateSNR(const Image &original, const Image &decoded);

// Helper to save/load EncodedData
void saveCompressed(const std::string &filename, const EncodedData &data);
EncodedData loadCompressed(const std::string &filename);
