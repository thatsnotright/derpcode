# Min Stack

Implement a stack of `int` that, in addition to the usual operations, returns its minimum element. **Every operation must run in O(1) time.**

## Methods to implement

```go
// NewMinStack returns an empty stack.
func NewMinStack() *MinStack

// Push adds x to the top of the stack.
func (s *MinStack) Push(x int)

// Pop removes the top element. Popping an empty stack is a no-op.
func (s *MinStack) Pop()

// Top returns the top element and true, or 0, false if the stack is empty.
func (s *MinStack) Top() (int, bool)

// GetMin returns the smallest element currently in the stack and true,
// or 0, false if the stack is empty.
func (s *MinStack) GetMin() (int, bool)
```

## Example

```go
s := NewMinStack()
s.Push(5)
s.Push(2)
s.Push(7)
s.GetMin() // 2, true
s.Pop()    // removes 7
s.Pop()    // removes 2
s.GetMin() // 5, true
s.Top()    // 5, true
```

## Constraints

- Values may be negative, zero, or duplicated; the minimum must stay correct when a duplicate minimum is popped.
- An empty stack reports `ok == false`; a stored `0` is distinct from "empty".
- Single-threaded use is fine.

## Follow-ups

1. Can you do it with a single stack (no parallel min stack)? Hint: store differences or `(value, minSoFar)` pairs.
2. Add `GetMax()` with the same guarantees.
3. How would you add `PopMin()` while keeping `Push`/`Pop` fast? What complexity can you hope for?
