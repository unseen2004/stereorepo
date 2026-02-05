#pragma once
#include <vector>
#include <string>
#include <cstdint>

struct Pixel {
    uint8_t b, g, r;
};

struct Image {
    int width;
    int height;
    std::vector<Pixel> pixels; // Row-major, top-left origin in memory.
};

// Reads an uncompressed RGB 24-bit TGA file.
Image readTGA(const std::string& filename);

// Writes an uncompressed RGB 24-bit TGA file.
void writeTGA(const std::string& filename, const Image& img);
