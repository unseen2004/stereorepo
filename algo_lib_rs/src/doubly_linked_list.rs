use crate::node_double::NodeDouble;

pub struct DoublyLinkedList<T> {
    head: Option<Box<NodeDouble<T>>>,
    tail: Option<*mut NodeDouble<T>>,
}

impl<T> DoublyLinkedList<T> {
    pub fn new() -> Self {
        DoublyLinkedList {
            head: None,
            tail: None,
        }
    }

    pub fn push_front(&mut self, val: T) {
        // Create a new boxed node where its `next` is the current head
        let mut new_node = Box::new(NodeDouble {
            val,
            next: self.head.take(),
            prev: None,
        });

        // Obtain a raw pointer to the new node's heap allocation.
        // This pointer does not take ownership; it is only used for linking `prev`/`tail`.
        let new_ptr: *mut NodeDouble<T> = &mut *new_node;

        // If there was an old head, update its `prev` to point to the new node.
        if let Some(ref mut old_head) = new_node.next {
            old_head.prev = Some(new_ptr);
        } else {
            // The list was empty, so the new node is also the tail.
            self.tail = Some(new_ptr);
        }

        // Make the new node the head of the list (it still owns the `next` chain).
        self.head = Some(new_node);
    }

    pub fn push_back(&mut self, val: T) {

    }
}