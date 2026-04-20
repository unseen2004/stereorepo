use anyhow::Result;
use clap::{Parser, Subcommand};
use std::fs;
use std::path::PathBuf;
use turboquant_compressor::{compress_file, decompress_file};
use walkdir::WalkDir;

#[derive(Parser)]
#[command(
    author,
    version,
    about = "TurboQuant Compressor - Efficient LLM Quantization Tool"
)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Compress {
        input: PathBuf,
        #[arg(short, long)]
        output: Option<PathBuf>,
        #[arg(short, long, default_value_t = 4)]
        bits: u8,
        #[arg(short, long, default_value_t = 42)]
        seed: u64,
    },
    Decompress {
        input: PathBuf,
        #[arg(short, long)]
        output: Option<PathBuf>,
    },
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    match cli.command {
        Commands::Compress {
            input,
            output,
            bits,
            seed,
        } => {
            if input.is_dir() {
                let out_dir = output.unwrap_or_else(|| input.with_extension("compressed"));
                fs::create_dir_all(&out_dir)?;
                for entry in WalkDir::new(&input).into_iter().filter_map(|e| e.ok()) {
                    if entry.path().extension().and_then(|s| s.to_str()) == Some("safetensors") {
                        let rel = entry.path().strip_prefix(&input)?;
                        let out_file = out_dir.join(rel).with_extension("tq.safetensors");
                        if let Some(parent) = out_file.parent() {
                            fs::create_dir_all(parent)?;
                        }
                        compress_file(entry.path(), &out_file, bits, seed)?;
                    }
                }
            } else {
                let out = output.unwrap_or_else(|| input.with_extension("tq.safetensors"));
                compress_file(&input, &out, bits, seed)?;
            }
            println!("Compression complete.");
        }
        Commands::Decompress { input, output } => {
            let out = output.unwrap_or_else(|| input.with_extension("dec.safetensors"));
            decompress_file(&input, &out)?;
            println!("Decompression complete.");
        }
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use bytemuck::cast_slice;
    use safetensors::{serialize, Dtype, View};

    struct TestTensor {
        data: Vec<u8>,
        shape: Vec<usize>,
    }
    impl View for TestTensor {
        fn dtype(&self) -> Dtype {
            Dtype::F32
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

    #[test]
    fn test_end_to_end() -> Result<()> {
        let temp_dir = tempfile::tempdir()?;
        let input_path = temp_dir.path().join("input.safetensors");
        let compressed_path = temp_dir.path().join("output.tq.safetensors");
        let decompressed_path = temp_dir.path().join("output.dec.safetensors");

        let dim = 64;
        let num_vecs = 2;
        let mut data = vec![0.5f32; dim * num_vecs];
        data[0] = 1.0;

        let tensors = vec![(
            "weight".to_string(),
            TestTensor {
                data: cast_slice(&data).to_vec(),
                shape: vec![num_vecs, dim],
            },
        )];

        let bytes = serialize(tensors, None).unwrap();
        fs::write(&input_path, bytes)?;

        compress_file(&input_path, &compressed_path, 4, 42)?;
        assert!(compressed_path.exists());

        decompress_file(&compressed_path, &decompressed_path)?;
        assert!(decompressed_path.exists());

        Ok(())
    }
}
