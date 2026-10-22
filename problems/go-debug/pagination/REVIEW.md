# Review notes: pagination

## Planted issues

1. **`listing.go: Page` off-by-one.** Pages are 1-based but `start := page * size` skips the first page (page 1 returns items 4-6). Fix: `(page - 1) * size`. Same family: **`TotalPages` truncates** (`10/3 == 3`, should be 4); use ceiling division `(n + size - 1) / size`, and guard `size < 1` (divide by zero panic). Caught by: `TestPageFirst`, `TestTotalPages`.
2. **`Page` panics instead of clamping/validating.** `end := start + size` is never clamped to `len(items)`, so the last (partial) page and any page past the end panic with "slice bounds out of range"; `page < 1` / `size < 1` give negative/zero bounds. Also `page*size` and `start+size` can overflow with hostile input. Fix: validate args, return an empty slice (or an error) past the end, compute `end` without overflow. Caught by: `TestPageLastIsPartial`, `TestPageBeyondEnd`, `TestPageInvalidArgs`.
3. **Internal slices escape: `Items()` returns `s.items`; `Page` returns a sub-slice of it.** Callers can modify catalogue elements, and because the sub-slice has spare capacity, `append(page, x)` overwrites the next page's first element in the shared backing array. Fix: return copies (`slices.Clone`), or at minimum a full slice expression `s.items[start:end:end]` (still lets callers mutate elements). Caught by: `TestCallerCannotMutateCatalogue`, `TestAppendToPageDoesNotCorruptNextPage`.
4. **`Search` filters in place: `matches := s.items[:0]`.** The "filter without allocating" idiom reuses the receiver's backing array, so every search compacts matches to the front of the catalogue and corrupts it (the "duplicated products in the CSV export" symptom). Fix: `matches := make([]Item, 0, ...)` or `var matches []Item`. Caught by: `TestSearchDoesNotChangeCatalogue` (needs names that alternate match/no-match to show the corruption).
5. **`NewService` stores the caller's slice.** Later caller writes change the service's data (and the service's in-place writes change the caller's). Fix: `slices.Clone` on the way in. Caught by: `TestServiceCopiesInput`.
6. **`Tag` writes to a nil map.** `tags` is never initialised in `NewService`, so the first `Tag` panics with "assignment to entry in nil map". (Reading a nil map is fine, writing is not.) Fix: `make(map[string][]int)` in the constructor. Caught by: `TestTagging`.
7. **Review only:** `ByTag` is O(tags x items) (build an ID index); `Search` returns `nil` vs empty slice inconsistently (matters for JSON `null` vs `[]`); no doc comment about concurrency (not safe for concurrent mutation); `Page` silently returns empty on invalid args - discuss returning an error or a `Page` struct with `Total`, `HasNext`.

## Coaching / best-practice talking points

- Slice internals: length vs capacity, shared backing arrays, `s[a:b:b]`, `slices.Clone`, "filter in place" is only safe when you own the slice.
- Defensive copies at API boundaries (constructor and getters); document when you do not copy for performance.
- Constructors initialise maps; zero-value usefulness vs constructor requirement.
- Boundary testing: first, last, partial last, exact fit, empty, beyond end, zero/negative/huge args; table-driven tests with `t.Run`.
- Off-by-one: define the contract (1-based pages) in the doc comment and in a test name.
- Avoid premature optimisation ("avoid allocations") that sacrifices correctness; measure first.
- Integer overflow in index arithmetic; panics vs errors in library code.
- `gofmt`, `go vet`, doc comments, naming (`size` vs `pageSize`).

## How to run the conversation

- Ask them to run the tests; the first failure is a panic. Ask what a "slice bounds out of range" message tells them (read the numbers).
- Ask "what does `append` do when there is spare capacity?" and have them predict a small example before running it.
- For the CSV duplication story, ask them to reproduce it in a test; they will land on the in-place filter themselves.
- Let them choose between copying and returning an error for invalid pages; ask what the UI needs.
