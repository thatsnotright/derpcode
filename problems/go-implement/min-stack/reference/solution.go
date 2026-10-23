package solution

type frame struct {
	val, min int
}

// MinStack is a stack that tracks its minimum element in O(1).
type MinStack struct {
	frames []frame // each frame remembers the minimum at that depth
}

// NewMinStack returns an empty stack.
func NewMinStack() *MinStack {
	return &MinStack{}
}

// Push adds x to the top of the stack.
func (s *MinStack) Push(x int) {
	m := x
	if n := len(s.frames); n > 0 && s.frames[n-1].min < m {
		m = s.frames[n-1].min
	}
	s.frames = append(s.frames, frame{val: x, min: m})
}

// Pop removes the top element. Popping an empty stack is a no-op.
func (s *MinStack) Pop() {
	if n := len(s.frames); n > 0 {
		s.frames = s.frames[:n-1]
	}
}

// Top returns the top element and true, or 0, false if the stack is empty.
func (s *MinStack) Top() (int, bool) {
	if len(s.frames) == 0 {
		return 0, false
	}
	return s.frames[len(s.frames)-1].val, true
}

// GetMin returns the smallest element currently in the stack and true,
// or 0, false if the stack is empty.
func (s *MinStack) GetMin() (int, bool) {
	if len(s.frames) == 0 {
		return 0, false
	}
	return s.frames[len(s.frames)-1].min, true
}
