# Meeting Room Scheduler

Implement a scheduler that books meetings into a fixed number of rooms, numbered `0` to `rooms-1`.

A meeting occupies the half-open interval `[start, end)` in integer minutes, so one meeting may start exactly when another ends in the same room. `Book` assigns the **lowest-numbered** room that is free for the whole interval.

## Methods to implement

```go
// Meeting is a booked time slot.
type Meeting struct {
	ID    string
	Start int
	End   int
}

func NewScheduler(rooms int) *Scheduler

// Book reserves the lowest-numbered free room for [start, end). It returns the
// room and true on success. It fails (0, false) if start >= end, if id is
// already booked, or if no room is free for the whole interval.
func (s *Scheduler) Book(id string, start, end int) (room int, ok bool)

// Cancel removes the booking with the given id. It reports whether the
// booking existed.
func (s *Scheduler) Cancel(id string) bool

// RoomOf returns the room a meeting is booked in.
func (s *Scheduler) RoomOf(id string) (room int, ok bool)

// Schedule returns the meetings in the room sorted by Start. It returns nil
// for an unknown room or an empty room.
func (s *Scheduler) Schedule(room int) []Meeting
```

## Example

```go
s := NewScheduler(2)
s.Book("a", 9, 10)    // 0, true
s.Book("b", 9, 11)    // 1, true  (room 0 is busy)
s.Book("c", 10, 12)   // 0, true  (back-to-back with "a")
s.Book("d", 10, 11)   // 0, false (both rooms busy)
s.Cancel("b")         // true
s.Book("d", 10, 11)   // 1, true
```

## Constraints

- `rooms >= 0`; with zero rooms every booking fails.
- Booking is greedy: it never moves an existing meeting.
- Single-threaded use is fine.

## Follow-ups

1. Make `Book` faster than scanning every meeting in every room. What data structure keeps each room's meetings ordered?
2. Offline version: given all meetings up front, what is the minimum number of rooms needed? (Sweep line or min-heap.)
3. Add `FreeSlots(room, from, to int) [][2]int` returning the gaps.
