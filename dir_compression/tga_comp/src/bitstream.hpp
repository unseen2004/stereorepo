#pragma once
#include <cstdint>
#include <iostream>
#include <vector>

class BitWriter {
  std::vector<uint8_t> &buffer;
  uint8_t currentByte;
  int bitCount;

public:
  BitWriter(std::vector<uint8_t> &buf)
      : buffer(buf), currentByte(0), bitCount(0) {}

  void writeBits(uint32_t val, int nBits) {
    for (int i = 0; i < nBits; ++i) {
      uint8_t bit = (val >> i) & 1; // Assuming LSB first for simplicity in
                                    // logic, but let's stick to standard order.
      // Actually, let's just write bits. Order: if I write bit 0 first.
      if (bit) {
        currentByte |= (1 << bitCount);
      }
      bitCount++;
      if (bitCount == 8) {
        buffer.push_back(currentByte);
        currentByte = 0;
        bitCount = 0;
      }
    }
  }

  void flush() {
    if (bitCount > 0) {
      buffer.push_back(currentByte);
      currentByte = 0;
      bitCount = 0; // Not explicitly needed but good for reset
    }
  }
};

class BitReader {
  const std::vector<uint8_t> &buffer;
  size_t byteIndex;
  int bitIndex;

public:
  BitReader(const std::vector<uint8_t> &buf)
      : buffer(buf), byteIndex(0), bitIndex(0) {}

  uint32_t readBits(int nBits) {
    uint32_t val = 0;
    for (int i = 0; i < nBits; ++i) {
      if (byteIndex >= buffer.size())
        return 0; // EOF safety

      uint8_t bit = (buffer[byteIndex] >> bitIndex) & 1;
      if (bit) {
        val |= (1 << i);
      }

      bitIndex++;
      if (bitIndex == 8) {
        bitIndex = 0;
        byteIndex++;
      }
    }
    return val;
  }
};
