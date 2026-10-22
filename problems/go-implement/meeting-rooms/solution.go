package solution

// Meeting is a booked time slot occupying [Start, End).
type Meeting struct {
	ID    string
	Start int
	End   int
}

// Scheduler books meetings into a fixed number of rooms.
type Scheduler struct {
	// TODO: add fields
}

// NewScheduler returns a scheduler with rooms rooms numbered 0..rooms-1.
func NewScheduler(rooms int) *Scheduler {
	// TODO: implement
	return &Scheduler{}
}

// Book reserves the lowest-numbered free room for [start, end). It returns the
// room and true on success. It fails (0, false) if start >= end, if id is
// already booked, or if no room is free for the whole interval.
func (s *Scheduler) Book(id string, start, end int) (room int, ok bool) {
	// TODO: implement
	return 0, false
}

// Cancel removes the booking with the given id. It reports whether the
// booking existed.
func (s *Scheduler) Cancel(id string) bool {
	// TODO: implement
	return false
}

// RoomOf returns the room a meeting is booked in.
func (s *Scheduler) RoomOf(id string) (room int, ok bool) {
	// TODO: implement
	return 0, false
}

// Schedule returns the meetings in the room sorted by Start. It returns nil
// for an unknown room or an empty room.
func (s *Scheduler) Schedule(room int) []Meeting {
	// TODO: implement
	return nil
}
