## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add paginated product list service

### What it does

Adds `listing.Service`, an in-memory catalogue used by the product list page:

- `Page(page, size)` returns one page of items (pages are 1-based).
- `TotalPages(size)` tells the UI how many pages exist.
- `Search(prefix)` returns the items whose name starts with a prefix.
- `Tag(id, tag)` / `ByTag(tag)` attach tags to items and look them up.
- `Items()` returns everything (used by the CSV export).

I tried to avoid allocations by reusing slices instead of copying where I could.

### How I tested it

Table of 10 fruit items; tests for the first page, the last page, total pages, search, and tagging.

The tests panic with "slice bounds out of range" on my machine and I can't tell if it's the tests or the
code. Also the CSV export sometimes shows duplicated products after someone uses the search box, which
I could not reproduce in a test.
