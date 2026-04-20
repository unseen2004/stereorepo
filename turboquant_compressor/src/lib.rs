use anyhow::{Context, Result};
use bytemuck::cast_slice;
use half::{bf16, f16};
use rayon::prelude::*;
use safetensors::{serialize, tensor::SafeTensors, Dtype, View};
use std::collections::HashMap;
use std::fs::File;
use std::io::Write;
use std::path::Path;
use turboquant::turboquant_mse::{QuantizedVector, TurboQuantMSE};

pub struct OutTensor {
    pub name: String,
    pub data: Vec<u8>,
    pub shape: Vec<usize>,
    pub dtype: Dtype,
}

impl View for OutTensor {
    fn dtype(&self) -> Dtype {
        self.dtype
    }
    fn shape(&self) -> &[usize] {
        &self.shape
    }
    fn data(&self) -> std::borrow::Cow<'_, [u8]> {
        std::borrow::Cow::Borrowed(&self.data)
    }
    fn data_len(&self) -> usize {
        self.data.len()
    }
}

pub fn compress_file(input: &Path, output: &Path, bits: u8, seed: u64) -> Result<()> {
    let file = File::open(input)?;
    let buffer = unsafe { memmap2::MmapOptions::new().map(&file)? };
    let tensors = SafeTensors::deserialize(&buffer)?;

    let mut out_tensors: Vec<(String, OutTensor)> = Vec::new();

    for (name, tensor) in tensors.tensors() {
        let shape = tensor.shape();
        if shape.len() < 2 {
            continue;
        }

        let dim = shape[shape.len() - 1];
        let num_vectors: usize = shape[..shape.len() - 1].iter().product();

        let f32_data: Vec<f32> = match tensor.dtype() {
            Dtype::F32 => cast_slice::<u8, f32>(tensor.data()).to_vec(),
            Dtype::F16 => cast_slice::<u8, f16>(tensor.data())
                .iter()
                .map(|f| f.to_f32())
                .collect(),
            Dtype::BF16 => cast_slice::<u8, bf16>(tensor.data())
                .iter()
                .map(|f| f.to_f32())
                .collect(),
            _ => continue,
        };

        let tq = TurboQuantMSE::new(dim, bits, seed)?;

        let results: Vec<(Vec<u8>, f32)> = f32_data
            .par_chunks_exact(dim)
            .map(|chunk| {
                let vec_f64: Vec<f64> = chunk.iter().map(|&x| x as f64).collect();
                let norm = (vec_f64.iter().map(|x| x * x).sum::<f64>()).sqrt();
                let normalized: Vec<f64> = vec_f64.iter().map(|&x| x / (norm + 1e-8)).collect();
                let q = tq.quantize(&normalized).unwrap();
                (q.indices, norm as f32)
            })
            .collect();

        let (quantized, norms): (Vec<u8>, Vec<f32>) = results.into_iter().fold(
            (Vec::new(), Vec::new()),
            |(mut q, mut n), (q_part, n_part)| {
                q.extend(q_part);
                n.push(n_part);
                (q, n)
            },
        );

        out_tensors.push((
            name.clone(),
            OutTensor {
                name: name.clone(),
                data: quantized,
                shape: shape.to_vec(),
                dtype: Dtype::U8,
            },
        ));
        out_tensors.push((
            format!("{}_norms", name),
            OutTensor {
                name: "".to_string(),
                data: cast_slice(&norms).to_vec(),
                shape: vec![num_vectors],
                dtype: Dtype::F32,
            },
        ));
    }

    let mut metadata = HashMap::new();
    metadata.insert("bits".to_string(), bits.to_string());
    metadata.insert("seed".to_string(), seed.to_string());

    let out_data = serialize(out_tensors, Some(metadata))?;
    File::create(output)?.write_all(&out_data)?;
    Ok(())
}

pub fn decompress_file(input: &Path, output: &Path) -> Result<()> {
    let file = File::open(input)?;
    let buffer = unsafe { memmap2::MmapOptions::new().map(&file)? };
    let (_, metadata_obj) = SafeTensors::read_metadata(&buffer)?;
    let metadata = metadata_obj.metadata();

    let (bits, seed) = if let Some(m) = metadata {
        let b = m
            .get("bits")
            .and_then(|s| s.parse::<u8>().ok())
            .unwrap_or(4);
        let s = m
            .get("seed")
            .and_then(|s| s.parse::<u64>().ok())
            .unwrap_or(42);
        (b, s)
    } else {
        (4, 42)
    };

    let tensors = SafeTensors::deserialize(&buffer)?;
    let mut out_tensors: Vec<(String, OutTensor)> = Vec::new();

    for (name, tensor) in tensors
        .tensors()
        .into_iter()
        .filter(|(n, _)| !n.ends_with("_norms"))
    {
        let norm_name = format!("{}_norms", name);
        let norms_tensor = tensors.tensor(&norm_name).context("Missing norms")?;
        let norms: &[f32] = cast_slice(norms_tensor.data());

        let shape = tensor.shape();
        let dim = shape[shape.len() - 1];
        let tq = TurboQuantMSE::new(dim, bits, seed)?;

        let q_data = tensor.data();
        let bytes_per_vec = q_data.len() / norms.len();

        let mut f32_flat = Vec::with_capacity(norms.len() * dim);
        for (i, chunk) in q_data.chunks_exact(bytes_per_vec).enumerate() {
            let q_vec = QuantizedVector {
                indices: chunk.to_vec(),
                bit_width: bits,
                dim,
            };
            let reconstructed = tq.dequantize(&q_vec).unwrap();
            for val in reconstructed {
                f32_flat.push((val * norms[i] as f64) as f32);
            }
        }

        out_tensors.push((
            name.clone(),
            OutTensor {
                name: name.clone(),
                data: cast_slice(&f32_flat).to_vec(),
                shape: shape.to_vec(),
                dtype: Dtype::F32,
            },
        ));
    }

    let out_data = serialize(out_tensors, None)?;
    File::create(output)?.write_all(&out_data)?;
    Ok(())
}
