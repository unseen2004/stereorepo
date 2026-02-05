#include "tga.hpp"
#include <cstring>
#include <fstream>
#include <iostream>
#include <stdexcept>

#pragma pack(push, 1)
struct TGAHeader {
  uint8_t idLength;
  uint8_t colorMapType;
  uint8_t imageType;
  uint16_t colorMapOrigin;
  uint16_t colorMapLength;
  uint8_t colorMapDepth;
  uint16_t xOrigin;
  uint16_t yOrigin;
  uint16_t width;
  uint16_t height;
  uint8_t pixelDepth;
  uint8_t imageDescriptor;
};
#pragma pack(pop)

Image readTGA(const std::string &filename) {
  std::ifstream file(filename, std::ios::binary);
  if (!file) {
    throw std::runtime_error("Cannot open file: " + filename);
  }

  TGAHeader header;
  file.read(reinterpret_cast<char *>(&header), sizeof(header));

  if (header.imageType != 2) { // 2 = Uncompressed RGB
    throw std::runtime_error(
        "Unsupported TGA type (only uncompressed RGB supported).");
  }
  if (header.pixelDepth != 24) {
    throw std::runtime_error(
        "Unsupported pixel depth (only 24-bit supported).");
  }

  // Skip ID
  file.seekg(header.idLength, std::ios::cur);

  Image img;
  img.width = header.width;
  img.height = header.height;
  img.pixels.resize(img.width * img.height);

  // Read raw pixel data
  std::vector<Pixel> raw(img.pixels.size());
  file.read(reinterpret_cast<char *>(raw.data()), raw.size() * 3);

  if (!file) {
    throw std::runtime_error("Error reading pixel data.");
  }

  bool origin_top = (header.imageDescriptor & 0x20) != 0;
  bool origin_right = (header.imageDescriptor & 0x10) != 0;

  // Normalize to top-left, left-to-right in memory.
  for (int y = 0; y < img.height; ++y) {
    int src_y = origin_top ? y : (img.height - 1 - y);
    for (int x = 0; x < img.width; ++x) {
      int src_x = origin_right ? (img.width - 1 - x) : x;
      img.pixels[y * img.width + x] =
          raw[src_y * img.width + src_x];
    }
  }

  return img;
}

void writeTGA(const std::string &filename, const Image &img) {
  std::ofstream file(filename, std::ios::binary);
  if (!file) {
    throw std::runtime_error("Cannot create file: " + filename);
  }

  TGAHeader header = {};
  header.imageType = 2; // Uncompressed RGB
  header.width = img.width;
  header.height = img.height;
  header.pixelDepth = 24;
  header.imageDescriptor = 0x20; // top-left origin, left-to-right

  file.write(reinterpret_cast<const char *>(&header), sizeof(header));
  file.write(reinterpret_cast<const char *>(img.pixels.data()),
             img.pixels.size() * 3);
}
