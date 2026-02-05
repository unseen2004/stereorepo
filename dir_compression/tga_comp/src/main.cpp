#include "compression.hpp"
#include "tga.hpp"
#include <iostream>
#include <string>
#include <vector>

void printUsage(const char *prog) {
  std::cout << "Usage:\n";
  std::cout << "  Encode: " << prog << " encode <input.tga> <output.bin> <k>\n";
  std::cout << "  Decode: " << prog << " decode <input.bin> <output.tga>\n";
  std::cout << "  Verify: " << prog << " verify <original.tga> <decoded.tga>\n";
  std::cout << "  Full Test: " << prog << " test <input.tga> <k>\n";
}

int main(int argc, char **argv) {
  if (argc < 2) {
    printUsage(argv[0]);
    return 1;
  }

  std::string mode = argv[1];

  try {
    if (mode == "encode") {
      if (argc != 5) {
        printUsage(argv[0]);
        return 1;
      }
      std::string inFile = argv[2];
      std::string outFile = argv[3];
      int k = std::stoi(argv[4]);
      if (k < 1 || k > 7) {
        std::cerr << "k must be between 1 and 7\n";
        return 1;
      }

      std::cout << "Reading " << inFile << "...\n";
      Image img = readTGA(inFile);

      std::cout << "Encoding with k=" << k << "...\n";
      EncodedData data = encode(img, k);

      std::cout << "Saving to " << outFile << "...\n";
      saveCompressed(outFile, data);

      std::cout << "Done.\n";

    } else if (mode == "decode") {
      if (argc != 4) {
        printUsage(argv[0]);
        return 1;
      }
      std::string inFile = argv[2];
      std::string outFile = argv[3];

      std::cout << "Loading " << inFile << "...\n";
      EncodedData data = loadCompressed(inFile);

      std::cout << "Decoding...\n";
      Image img = decode(data);

      std::cout << "Saving to " << outFile << "...\n";
      writeTGA(outFile, img);

      std::cout << "Done.\n";

    } else if (mode == "verify") {
      if (argc != 4) {
        printUsage(argv[0]);
        return 1;
      }
      std::string origFile = argv[2];
      std::string decFile = argv[3];

      Image orig = readTGA(origFile);
      Image dec = readTGA(decFile);

      double mse = calculateMSE(orig, dec);
      double snr = calculateSNR(orig, dec);

      std::cout << "MSE: " << mse << "\n";
      std::cout << "SNR: " << snr << " dB\n";

      // Also component-wise
      // Re-calc for specific colors or just overall is fine for now.
      // Requirement says "error for whole image and individual color
      // components".

      double mse_r = 0, mse_g = 0, mse_b = 0;
      for (size_t i = 0; i < orig.pixels.size(); ++i) {
        double d;
        d = orig.pixels[i].r - dec.pixels[i].r;
        mse_r += d * d;
        d = orig.pixels[i].g - dec.pixels[i].g;
        mse_g += d * d;
        d = orig.pixels[i].b - dec.pixels[i].b;
        mse_b += d * d;
      }
      mse_r /= orig.pixels.size();
      mse_g /= orig.pixels.size();
      mse_b /= orig.pixels.size();

      std::cout << "MSE (R): " << mse_r << "\n";
      std::cout << "MSE (G): " << mse_g << "\n";
      std::cout << "MSE (B): " << mse_b << "\n";

    } else if (mode == "test") {
      if (argc != 4) {
        printUsage(argv[0]);
        return 1;
      }
      std::string inFile = argv[2];
      int k = std::stoi(argv[3]);

      std::cout << "Reading " << inFile << "...\n";
      Image img = readTGA(inFile);

      std::cout << "Encoding k=" << k << "...\n";
      EncodedData data = encode(img, k);

      std::cout << "Decoding...\n";
      Image dec = decode(data);

      double mse = calculateMSE(img, dec);
      double snr = calculateSNR(img, dec);

      std::cout << "MSE: " << mse << "\n";
      std::cout << "SNR: " << snr << " dB\n";

    } else {
      printUsage(argv[0]);
      return 1;
    }
  } catch (const std::exception &e) {
    std::cerr << "Error: " << e.what() << "\n";
    return 1;
  }

  return 0;
}
