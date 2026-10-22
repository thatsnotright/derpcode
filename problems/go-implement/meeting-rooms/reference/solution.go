package solution

import "sort"

// Meeting is a booked time slot occupying [Start, End).
type Meeting struct {
	ID    string
	Start int
	End   int
}

// Scheduler books meetings into a fixed number of rooms.
type Scheduler struct {
	rooms  [][]Meeting // each room's meetings sorted by Start, non-overlapping
	byID   map[string]int
	nRooms int
}

// NewScheduler returns a scheduler with rooms rooms numbered 0..rooms-1.
func NewScheduler(rooms int) *Scheduler {
	if rooms < 0 {
		rooms = 0
	}
	return &Scheduler{
		rooms:  make([][]Meeting, rooms),
		byID:   make(map[string]int),
		nRooms: rooms,
	}
}

// Book reserves the lowest-numbered free room for [start, end). It returns the
// room and true on success. It fails (0, false) if start >= end, if id is
// already booked, or if no room is free for the whole interval.
func (s *Scheduler) Book(id string, start, end int) (room int, ok bool) {
	if start >= end {
		return 0, false
	}
	if _, dup := s.byID[id]; dup {
		return 0, false
	}
	for r, meetings := range s.rooms {
		// First meeting that ends after start; the slot is free iff it
		// also begins at or after end.
		i := sort.Search(len(meetings), func(i int) bool { return meetings[i].End > start })
		if i < len(meetings) && meetings[i].Start < end {
			continue
		}
		m := Meeting{ID: id, Start: start, End: end}
		meetings = append(meetings, Meeting{})
		copy(meetings[i+1:], meetings[i:])
		meetings[i] = m
		s.rooms[r] = meetings
		s.byID[id] = r
		return r, true
	}
	return 0, false
}

// Cancel removes the booking with the given id. It reports whether the
// booking existed.
func (s *Scheduler) Cancel(id string) bool {
	r, ok := s.byID[id]
	if !ok {
		return false
	}
	delete(s.byID, id)
	meetings := s.rooms[r]
	for i, m := range meetings {
		if m.ID == id {
			s.rooms[r] = append(meetings[:i], meetings[i+1:]...)
			break
		}
	}
	return true
}

// RoomOf returns the room a meeting is booked in.
func (s *Scheduler) RoomOf(id string) (room int, ok bool) {
	room, ok = s.byID[id]
	return room, ok
}

// Schedule returns the meetings in the room sorted by Start. It returns nil
// for an unknown room or an empty room.
func (s *Scheduler) Schedule(room int) []Meeting {
	if room < 0 || room >= s.nRooms || len(s.rooms[room]) == 0 {
		return nil
	}
	out := make([]Meeting, len(s.rooms[room]))
	copy(out, s.rooms[room])
	return out
}
