package solution

import (
	"reflect"
	"testing"
)

func TestBookLowestFreeRoom(t *testing.T) {
	s := NewScheduler(3)
	tests := []struct {
		id         string
		start, end int
		wantRoom   int
		wantOK     bool
	}{
		{"a", 9, 11, 0, true},
		{"b", 10, 12, 1, true},
		{"c", 10, 11, 2, true},
		{"d", 10, 11, 0, false},
		{"e", 11, 12, 0, true}, // room 0 frees at 11
	}
	for _, tt := range tests {
		room, ok := s.Book(tt.id, tt.start, tt.end)
		if room != tt.wantRoom || ok != tt.wantOK {
			t.Errorf("Book(%q, %d, %d) = (%d, %v), want (%d, %v)",
				tt.id, tt.start, tt.end, room, ok, tt.wantRoom, tt.wantOK)
		}
	}
}

func TestBackToBackIsAllowed(t *testing.T) {
	s := NewScheduler(1)
	if _, ok := s.Book("a", 9, 10); !ok {
		t.Fatal("first booking should succeed")
	}
	if _, ok := s.Book("b", 10, 11); !ok {
		t.Error("meeting starting at the previous end should be allowed")
	}
	if _, ok := s.Book("c", 8, 9); !ok {
		t.Error("meeting ending at the next start should be allowed")
	}
}

func TestOverlapVariants(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
	}{
		{"identical", 10, 20},
		{"inside", 12, 18},
		{"contains", 5, 25},
		{"overlaps start", 5, 15},
		{"overlaps end", 15, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScheduler(1)
			s.Book("base", 10, 20)
			if room, ok := s.Book("new", tt.start, tt.end); ok {
				t.Errorf("Book overlapping meeting succeeded in room %d", room)
			}
		})
	}
}

func TestInvalidIntervals(t *testing.T) {
	s := NewScheduler(2)
	for _, iv := range [][2]int{{5, 5}, {6, 5}} {
		if _, ok := s.Book("x", iv[0], iv[1]); ok {
			t.Errorf("Book(%d, %d) should fail", iv[0], iv[1])
		}
	}
	if _, ok := s.RoomOf("x"); ok {
		t.Error("failed booking must not be recorded")
	}
}

func TestDuplicateIDRejected(t *testing.T) {
	s := NewScheduler(2)
	s.Book("a", 1, 2)
	if _, ok := s.Book("a", 5, 6); ok {
		t.Error("duplicate id should be rejected")
	}
	if got := s.Schedule(0); len(got) != 1 {
		t.Errorf("Schedule(0) = %v, want one meeting", got)
	}
}

func TestZeroRooms(t *testing.T) {
	s := NewScheduler(0)
	if _, ok := s.Book("a", 1, 2); ok {
		t.Error("booking with zero rooms should fail")
	}
}

func TestCancelFreesRoom(t *testing.T) {
	s := NewScheduler(2)
	s.Book("a", 9, 10)
	s.Book("b", 9, 10)
	if !s.Cancel("a") {
		t.Fatal("Cancel(a) should succeed")
	}
	if s.Cancel("a") {
		t.Error("second Cancel(a) should report false")
	}
	if s.Cancel("nope") {
		t.Error("Cancel of unknown id should report false")
	}
	if room, ok := s.Book("c", 9, 10); !ok || room != 0 {
		t.Errorf("Book(c) = (%d, %v), want (0, true)", room, ok)
	}
	if _, ok := s.RoomOf("a"); ok {
		t.Error("cancelled meeting should have no room")
	}
}

func TestRoomOf(t *testing.T) {
	s := NewScheduler(2)
	s.Book("a", 9, 10)
	s.Book("b", 9, 10)
	for id, want := range map[string]int{"a": 0, "b": 1} {
		if room, ok := s.RoomOf(id); !ok || room != want {
			t.Errorf("RoomOf(%q) = (%d, %v), want (%d, true)", id, room, ok, want)
		}
	}
	if _, ok := s.RoomOf("zzz"); ok {
		t.Error("RoomOf unknown id should be false")
	}
}

func TestScheduleSortedByStart(t *testing.T) {
	s := NewScheduler(1)
	s.Book("late", 30, 40)
	s.Book("early", 0, 10)
	s.Book("mid", 10, 20)
	want := []Meeting{
		{"early", 0, 10},
		{"mid", 10, 20},
		{"late", 30, 40},
	}
	if got := s.Schedule(0); !reflect.DeepEqual(got, want) {
		t.Errorf("Schedule(0) = %v, want %v", got, want)
	}
}

func TestScheduleUnknownOrEmptyRoom(t *testing.T) {
	s := NewScheduler(2)
	s.Book("a", 1, 2)
	if got := s.Schedule(1); got != nil {
		t.Errorf("Schedule(empty room) = %v, want nil", got)
	}
	for _, r := range []int{-1, 2, 99} {
		if got := s.Schedule(r); got != nil {
			t.Errorf("Schedule(%d) = %v, want nil", r, got)
		}
	}
}
