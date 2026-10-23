package solution

import "testing"

func checkTop(t *testing.T, s *MinStack, want int, wantOK bool) {
	t.Helper()
	if got, ok := s.Top(); got != want || ok != wantOK {
		t.Errorf("Top() = (%d, %v), want (%d, %v)", got, ok, want, wantOK)
	}
}

func checkMin(t *testing.T, s *MinStack, want int, wantOK bool) {
	t.Helper()
	if got, ok := s.GetMin(); got != want || ok != wantOK {
		t.Errorf("GetMin() = (%d, %v), want (%d, %v)", got, ok, want, wantOK)
	}
}

func TestEmptyStack(t *testing.T) {
	s := NewMinStack()
	checkTop(t, s, 0, false)
	checkMin(t, s, 0, false)
	s.Pop() // must not panic
	checkTop(t, s, 0, false)
}

func TestPushTopMin(t *testing.T) {
	s := NewMinStack()
	s.Push(5)
	checkTop(t, s, 5, true)
	checkMin(t, s, 5, true)
	s.Push(2)
	checkTop(t, s, 2, true)
	checkMin(t, s, 2, true)
	s.Push(7)
	checkTop(t, s, 7, true)
	checkMin(t, s, 2, true)
}

func TestMinRestoredAfterPop(t *testing.T) {
	s := NewMinStack()
	for _, x := range []int{5, 2, 7} {
		s.Push(x)
	}
	s.Pop() // 7
	checkMin(t, s, 2, true)
	s.Pop() // 2
	checkMin(t, s, 5, true)
	checkTop(t, s, 5, true)
}

func TestDuplicateMinimum(t *testing.T) {
	s := NewMinStack()
	s.Push(3)
	s.Push(1)
	s.Push(1)
	s.Pop()
	checkMin(t, s, 1, true)
	s.Pop()
	checkMin(t, s, 3, true)
}

func TestNegativeAndZero(t *testing.T) {
	s := NewMinStack()
	s.Push(0)
	checkMin(t, s, 0, true)
	checkTop(t, s, 0, true)
	s.Push(-4)
	checkMin(t, s, -4, true)
	s.Push(-1)
	checkMin(t, s, -4, true)
	s.Pop()
	s.Pop()
	checkMin(t, s, 0, true)
}

func TestDescendingThenAscending(t *testing.T) {
	s := NewMinStack()
	for _, x := range []int{9, 8, 7, 6} {
		s.Push(x)
		checkMin(t, s, x, true)
	}
	for _, want := range []int{7, 8, 9} {
		s.Pop()
		checkMin(t, s, want, true)
	}
}

func TestDrainThenReuse(t *testing.T) {
	s := NewMinStack()
	s.Push(1)
	s.Pop()
	s.Pop()
	checkMin(t, s, 0, false)
	s.Push(10)
	checkMin(t, s, 10, true)
	checkTop(t, s, 10, true)
}

func TestLargeSequence(t *testing.T) {
	s := NewMinStack()
	const n = 10000
	for i := n; i > 0; i-- {
		s.Push(i)
	}
	checkMin(t, s, 1, true)
	for i := 1; i < n; i++ {
		s.Pop()
		checkMin(t, s, i+1, true)
	}
}
