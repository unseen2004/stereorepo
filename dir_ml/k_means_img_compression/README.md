# K-Means Image Compression

Image compression using K-Means clustering algorithm to reduce color palette.

---

## Overview

Applies K-Means clustering to compress images by reducing the number of colors. Each pixel is assigned to the nearest cluster centroid, effectively quantizing the color space.

---

## Requirements

Install Jupyter and required packages:
```bash
pip install jupyter numpy matplotlib pillow scikit-learn
```

---

## Running

Launch the Jupyter notebook:
```bash
jupyter notebook k_mean.ipynb
```

---

## Features

- Image color quantization
- Configurable compression levels (K values)
- Visual comparison of original vs compressed images
- File size reduction analysis

---

## How It Works

1. Load image data
2. Apply K-Means to group similar colors
3. Replace each pixel color with its cluster centroid
4. Display compressed result and metrics
