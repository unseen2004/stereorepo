use pyo3::prelude::*;
use turboquant::turboquant_mse::TurboQuantMSE;

#[pyfunction]
#[pyo3(signature = (data, dim, bits=4))]
fn compress(data: Vec<f32>, dim: usize, bits: u32) -> PyResult<(Vec<u8>, Vec<f32>)> {
    let tq = TurboQuantMSE::new(dim, bits as u8, 42)
        .map_err(|e| pyo3::exceptions::PyValueError::new_err(e.to_string()))?;
    
    let mut quantized = Vec::new();
    let mut norms = Vec::new();
    
    for chunk in data.chunks_exact(dim) {
        let vec_f64: Vec<f64> = chunk.iter().map(|&x| x as f64).collect();
        let norm = (vec_f64.iter().map(|x| x * x).sum::<f64>()).sqrt();
        let normalized: Vec<f64> = vec_f64.iter().map(|&x| x / (norm + 1e-8)).collect();
        
        norms.push(norm as f32);
        let q = tq.quantize(&normalized)
            .map_err(|e| pyo3::exceptions::PyValueError::new_err(e.to_string()))?;
            
        quantized.extend(q.indices);
    }
    
    Ok((quantized, norms))
}

#[pymodule]
fn turboquant_compressor(_py: Python, m: &PyModule) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(compress, m)?)?;
    Ok(())
}
