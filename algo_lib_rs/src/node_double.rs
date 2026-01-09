#[derive(Debug)]
pub struct NodeDouble<T> {
    pub val: T,
    pub next: Option<Box<NodeDouble<T>>>,
    pub prev: Option<*mut NodeDouble<T>>,
}

impl<T> NodeDouble<T> {
    pub fn new(val: T) -> Self {
        NodeDouble {
            val,
            next: None,
            prev: None,
        }
    }
}
