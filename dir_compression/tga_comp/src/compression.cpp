#include "compression.hpp"
#include <algorithm>
#include <cmath>
#include <fstream>
#include <iostream>
#include <limits>

// --- Helper Functions ---

static double clamp(double v) {
  if (v < 0)
    return 0;
  if (v > 255)
    return 255;
  return v;
}

// Haar-like 1D Decomposition: L=(a+b)/2, H=(a-b)/2
static void decomposeRow(const std::vector<double> &input,
                         std::vector<double> &L, std::vector<double> &H) {
  // Process pairs; caller provides even-length input (pad if needed).
  size_t width = input.size();
  L.clear();
  H.clear();
  L.reserve(width / 2);
  H.reserve(width / 2);

  for (size_t i = 0; i + 1 < width; i += 2) {
    double a = input[i];
    double b = input[i + 1];
    L.push_back((a + b) / 2.0);
    H.push_back((a - b) / 2.0);
  }
}

static void reconstructRow(const std::vector<double> &L,
                           const std::vector<double> &H,
                           std::vector<double> &output) {
  output.clear();
  output.reserve(L.size() * 2);
  for (size_t i = 0; i < L.size(); ++i) {
    double l = L[i];
    double h = H[i];
    // a = l + h
    // b = l - h
    output.push_back(clamp(l + h));
    output.push_back(clamp(l - h));
  }
}

static void decompose2D(const std::vector<double> &input, int width, int height,
                        std::vector<double> &LL, std::vector<double> &HL,
                        std::vector<double> &LH, std::vector<double> &HH) {
  int half_w = width / 2;
  int half_h = height / 2;
  std::vector<double> Lx(height * half_w);
  std::vector<double> Hx(height * half_w);

  // Row-wise decomposition.
  for (int y = 0; y < height; ++y) {
    std::vector<double> row(width);
    for (int x = 0; x < width; ++x) {
      row[x] = input[y * width + x];
    }
    std::vector<double> L, H;
    decomposeRow(row, L, H);
    for (int x = 0; x < half_w; ++x) {
      Lx[y * half_w + x] = L[x];
      Hx[y * half_w + x] = H[x];
    }
  }

  LL.assign(half_w * half_h, 0.0);
  HL.assign(half_w * half_h, 0.0);
  LH.assign(half_w * half_h, 0.0);
  HH.assign(half_w * half_h, 0.0);

  // Column-wise decomposition.
  for (int x = 0; x < half_w; ++x) {
    std::vector<double> colL(height), colH(height);
    for (int y = 0; y < height; ++y) {
      colL[y] = Lx[y * half_w + x];
      colH[y] = Hx[y * half_w + x];
    }
    std::vector<double> Lc, Hc;
    decomposeRow(colL, Lc, Hc);
    for (int y = 0; y < half_h; ++y) {
      LL[y * half_w + x] = Lc[y];
      HL[y * half_w + x] = Hc[y];
    }

    decomposeRow(colH, Lc, Hc);
    for (int y = 0; y < half_h; ++y) {
      LH[y * half_w + x] = Lc[y];
      HH[y * half_w + x] = Hc[y];
    }
  }
}

static void reconstruct2D(const std::vector<double> &LL,
                          const std::vector<double> &HL,
                          const std::vector<double> &LH,
                          const std::vector<double> &HH, int width, int height,
                          std::vector<double> &output) {
  int half_w = width / 2;
  int half_h = height / 2;
  std::vector<double> Lx(height * half_w);
  std::vector<double> Hx(height * half_w);

  // Column-wise reconstruction.
  for (int x = 0; x < half_w; ++x) {
    std::vector<double> Lc(half_h), Hc(half_h), col(height);
    for (int y = 0; y < half_h; ++y) {
      Lc[y] = LL[y * half_w + x];
      Hc[y] = HL[y * half_w + x];
    }
    reconstructRow(Lc, Hc, col);
    for (int y = 0; y < height; ++y) {
      Lx[y * half_w + x] = col[y];
    }

    for (int y = 0; y < half_h; ++y) {
      Lc[y] = LH[y * half_w + x];
      Hc[y] = HH[y * half_w + x];
    }
    reconstructRow(Lc, Hc, col);
    for (int y = 0; y < height; ++y) {
      Hx[y * half_w + x] = col[y];
    }
  }

  // Row-wise reconstruction.
  output.assign(width * height, 0.0);
  for (int y = 0; y < height; ++y) {
    std::vector<double> rowL(half_w), rowH(half_w), rowOut;
    for (int x = 0; x < half_w; ++x) {
      rowL[x] = Lx[y * half_w + x];
      rowH[x] = Hx[y * half_w + x];
    }
    reconstructRow(rowL, rowH, rowOut);
    for (int x = 0; x < width; ++x) {
      output[y * width + x] = rowOut[x];
    }
  }
}

// DPCM
static std::vector<double> encodeDPCM(const std::vector<double> &in) {
  std::vector<double> out;
  out.reserve(in.size());
  if (in.empty())
    return out;

  out.push_back(in[0]);
  for (size_t i = 1; i < in.size(); ++i) {
    out.push_back(in[i] - in[i - 1]);
  }
  return out;
}

static std::vector<double> decodeDPCM(const std::vector<double> &in) {
  std::vector<double> out;
  out.reserve(in.size());
  if (in.empty())
    return out;

  double current = in[0];
  out.push_back(current);
  for (size_t i = 1; i < in.size(); ++i) {
    current = current + in[i];
    out.push_back(current);
  }
  return out;
}

// --- Quantizer ---

void Quantizer::train(std::vector<double> &data, int k) {
  int num_levels = 1 << k;
  reconstruction_values.resize(num_levels);

  if (data.empty()) {
    std::fill(reconstruction_values.begin(), reconstruction_values.end(), 0.0);
    return;
  }

  std::sort(data.begin(), data.end());

  // Percentile method: Split into equal chunks.
  // Chunk size
  size_t n = data.size();
  for (int i = 0; i < num_levels; ++i) {
    size_t start_idx = (size_t)((double)i * n / num_levels);
    size_t end_idx = (size_t)((double)(i + 1) * n / num_levels);
    if (end_idx > n)
      end_idx = n;
    if (start_idx >= end_idx) {
      // Should not happen if n large enough, else take nearest
      if (start_idx > 0)
        reconstruction_values[i] = data[start_idx - 1];
      else
        reconstruction_values[i] = 0;
      continue;
    }

    // Average of the chunk
    double sum = 0;
    for (size_t j = start_idx; j < end_idx; ++j) {
      sum += data[j];
    }
    reconstruction_values[i] = sum / (end_idx - start_idx);
  }
}

uint32_t Quantizer::quantize(double val) const {
  // Find nearest reconstruction value
  double min_dist = std::numeric_limits<double>::max();
  uint32_t best_idx = 0;

  // Since reconstruction_values are sorted (because data was sorted and we took
  // averages of consecutive chunks), we could binary search. But linear is fine
  // for k <= 7 (max 128 items).
  for (size_t i = 0; i < reconstruction_values.size(); ++i) {
    double dist = std::abs(val - reconstruction_values[i]);
    if (dist < min_dist) {
      min_dist = dist;
      best_idx = i;
    }
  }
  return best_idx;
}

double Quantizer::dequantize(uint32_t index) const {
  if (index >= reconstruction_values.size())
    return 0;
  return reconstruction_values[index];
}

// --- Pipeline ---

EncodedData encode(const Image &img, int k) {
  EncodedData enc;
  enc.width = img.width;
  enc.height = img.height;
  enc.k = k;

  std::vector<double> R, G, B;
  size_t num_pixels = img.pixels.size();
  R.reserve(num_pixels);
  G.reserve(num_pixels);
  B.reserve(num_pixels);

  for (const auto &p : img.pixels) {
    R.push_back(p.r);
    G.push_back(p.g);
    B.push_back(p.b);
  }

  std::vector<double> *channels[3] = {&R, &G, &B};

  // --- Pass 1: Statistics Collection (Open Loop Diffs ONLY) ---
  std::vector<double> train_L[3];
  std::vector<double> train_H[3];

  for (int c = 0; c < 3; ++c) {
    const auto &origin = *channels[c];
    int w = enc.width;
    int h = enc.height;
    int padded_w = (w % 2 == 0) ? w : (w + 1);
    int padded_h = (h % 2 == 0) ? h : (h + 1);

    std::vector<double> padded(padded_w * padded_h);
    for (int y = 0; y < padded_h; ++y) {
      int src_y = (y < h) ? y : (h - 1);
      for (int x = 0; x < padded_w; ++x) {
        int src_x = (x < w) ? x : (w - 1);
        padded[y * padded_w + x] = origin[src_y * w + src_x];
      }
    }

    std::vector<double> LL, HL, LH, HH;
    decompose2D(padded, padded_w, padded_h, LL, HL, LH, HH);

    // Collect H (all high bands).
    train_H[c].insert(train_H[c].end(), HL.begin(), HL.end());
    train_H[c].insert(train_H[c].end(), LH.begin(), LH.end());
    train_H[c].insert(train_H[c].end(), HH.begin(), HH.end());

    // Collect LL diffs (i=1..N), ignore LL[0] for stats.
    if (LL.size() > 1) {
      for (size_t i = 1; i < LL.size(); ++i) {
        train_L[c].push_back(LL[i] - LL[i - 1]);
      }
    }
  }

  // Train Quantizers
  for (int c = 0; c < 3; ++c) {
    enc.quant_L[c].train(train_L[c], k);
    enc.quant_H[c].train(train_H[c], k);
  }

  // --- Pass 2: Encoding (Closed Loop) ---
  BitWriter writer(enc.bitstream);

  for (int c = 0; c < 3; ++c) {
    const auto &origin = *channels[c];
    int w = enc.width;
    int h = enc.height;
    int padded_w = (w % 2 == 0) ? w : (w + 1);
    int padded_h = (h % 2 == 0) ? h : (h + 1);

    std::vector<double> padded(padded_w * padded_h);
    for (int y = 0; y < padded_h; ++y) {
      int src_y = (y < h) ? y : (h - 1);
      for (int x = 0; x < padded_w; ++x) {
        int src_x = (x < w) ? x : (w - 1);
        padded[y * padded_w + x] = origin[src_y * w + src_x];
      }
    }

    std::vector<double> LL, HL, LH, HH;
    decompose2D(padded, padded_w, padded_h, LL, HL, LH, HH);

    // Encode LL
    // LL[0] stored as fixed-point (Q8.8) to preserve fractional averages.
    double pred = 0;
    if (!LL.empty()) {
      uint16_t start_val =
          (uint16_t)std::lround(clamp(LL[0]) * 256.0);
      writer.writeBits(start_val, 16);
      pred = (double)start_val / 256.0;
    }

    // LL[1..] Closed Loop DPCM
    for (size_t i = 1; i < LL.size(); ++i) {
      double val = LL[i];
      double diff = val - pred;

      uint32_t q = enc.quant_L[c].quantize(diff);
      double rec_diff = enc.quant_L[c].dequantize(q);

      pred = pred + rec_diff;
      writer.writeBits(q, k);
    }

    // Encode high bands (HL, LH, HH).
    for (double val : HL) {
      uint32_t q = enc.quant_H[c].quantize(val);
      writer.writeBits(q, k);
    }
    for (double val : LH) {
      uint32_t q = enc.quant_H[c].quantize(val);
      writer.writeBits(q, k);
    }
    for (double val : HH) {
      uint32_t q = enc.quant_H[c].quantize(val);
      writer.writeBits(q, k);
    }
  }
  writer.flush();

  return enc;
}

Image decode(const EncodedData &enc) {
  Image img;
  img.width = enc.width;
  img.height = enc.height;

  // We need to reconstruct channels.
  int w = enc.width;
  int h = enc.height;
  int padded_w = (w % 2 == 0) ? w : (w + 1);
  int padded_h = (h % 2 == 0) ? h : (h + 1);
  int band_w = padded_w / 2;
  int band_h = padded_h / 2;

  std::vector<double> reconstructed_channels[3];
  BitReader reader(enc.bitstream);

  for (int c = 0; c < 3; ++c) {
    std::vector<double> channel_pixels;
    channel_pixels.reserve(enc.width * enc.height);

    std::vector<double> LL(band_w * band_h);
    std::vector<double> HL(band_w * band_h);
    std::vector<double> LH(band_w * band_h);
    std::vector<double> HH(band_w * band_h);

    // Read LL
    double pred = 0;
    if (!LL.empty()) {
      uint32_t start_val = reader.readBits(16);
      pred = (double)start_val / 256.0;
      LL[0] = pred;
    }
    for (size_t i = 1; i < LL.size(); ++i) {
      uint32_t q = reader.readBits(enc.k);
      double diff = enc.quant_L[c].dequantize(q);
      pred += diff;
      LL[i] = pred;
    }

    // Read high bands.
    for (size_t i = 0; i < HL.size(); ++i) {
      uint32_t q = reader.readBits(enc.k);
      HL[i] = enc.quant_H[c].dequantize(q);
    }
    for (size_t i = 0; i < LH.size(); ++i) {
      uint32_t q = reader.readBits(enc.k);
      LH[i] = enc.quant_H[c].dequantize(q);
    }
    for (size_t i = 0; i < HH.size(); ++i) {
      uint32_t q = reader.readBits(enc.k);
      HH[i] = enc.quant_H[c].dequantize(q);
    }

    std::vector<double> padded;
    reconstruct2D(LL, HL, LH, HH, padded_w, padded_h, padded);

    for (int y = 0; y < h; ++y) {
      for (int x = 0; x < w; ++x) {
        channel_pixels.push_back(padded[y * padded_w + x]);
      }
    }
    reconstructed_channels[c] = channel_pixels;
  }

  img.pixels.resize(enc.width * enc.height);
  for (size_t i = 0; i < img.pixels.size(); ++i) {
    img.pixels[i].r = (uint8_t)clamp(reconstructed_channels[0][i]);
    img.pixels[i].g = (uint8_t)clamp(reconstructed_channels[1][i]);
    img.pixels[i].b = (uint8_t)clamp(reconstructed_channels[2][i]);
  }

  return img;
}

void saveCompressed(const std::string &filename, const EncodedData &data) {
  std::ofstream file(filename, std::ios::binary);
  file.write((char *)&data.width, sizeof(data.width));
  file.write((char *)&data.height, sizeof(data.height));
  file.write((char *)&data.k, sizeof(data.k));

  int num_levels = 1 << data.k;
  for (int c = 0; c < 3; ++c) {
    file.write((char *)data.quant_L[c].reconstruction_values.data(),
               num_levels * sizeof(double));
  }
  for (int c = 0; c < 3; ++c) {
    file.write((char *)data.quant_H[c].reconstruction_values.data(),
               num_levels * sizeof(double));
  }

  // Write bitstream size then data for safety? Or just EOF.
  // Let's write size.
  uint32_t bsize = data.bitstream.size();
  file.write((char *)&bsize, sizeof(bsize));
  file.write((char *)data.bitstream.data(), bsize);
}

EncodedData loadCompressed(const std::string &filename) {
  EncodedData data;
  std::ifstream file(filename, std::ios::binary);
  if (!file)
    throw std::runtime_error("Cannot open file");

  file.read((char *)&data.width, sizeof(data.width));
  file.read((char *)&data.height, sizeof(data.height));
  file.read((char *)&data.k, sizeof(data.k));

  int num_levels = 1 << data.k;
  for (int c = 0; c < 3; ++c) {
    data.quant_L[c].reconstruction_values.resize(num_levels);
    file.read((char *)data.quant_L[c].reconstruction_values.data(),
              num_levels * sizeof(double));
  }
  for (int c = 0; c < 3; ++c) {
    data.quant_H[c].reconstruction_values.resize(num_levels);
    file.read((char *)data.quant_H[c].reconstruction_values.data(),
              num_levels * sizeof(double));
  }

  uint32_t bsize;
  file.read((char *)&bsize, sizeof(bsize));
  data.bitstream.resize(bsize);
  file.read((char *)data.bitstream.data(), bsize);

  return data;
}

double calculateMSE(const Image &original, const Image &decoded) {
  if (original.pixels.size() != decoded.pixels.size())
    return -1.0;

  double error = 0;
  for (size_t i = 0; i < original.pixels.size(); ++i) {
    double dr = (double)original.pixels[i].r - decoded.pixels[i].r;
    double dg = (double)original.pixels[i].g - decoded.pixels[i].g;
    double db = (double)original.pixels[i].b - decoded.pixels[i].b;
    error += dr * dr + dg * dg + db * db;
  }
  return error / (original.pixels.size() * 3.0); // Per sample MSE
}

double calculateSNR(const Image &original, const Image &decoded) {
  double mse = calculateMSE(original, decoded);
  if (mse == 0)
    return 99999.0; // Perfect

  double signal = 0;
  for (size_t i = 0; i < original.pixels.size(); ++i) {
    double r = original.pixels[i].r;
    double g = original.pixels[i].g;
    double b = original.pixels[i].b;
    signal += r * r + g * g + b * b;
  }
  double power = signal / (original.pixels.size() * 3.0);

  return 10.0 * std::log10(power / mse);
}
