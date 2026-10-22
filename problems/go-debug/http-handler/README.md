## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add JSON API for todos

### What it does

Adds a tiny `net/http` JSON API backed by an in-memory store:

- `POST /todos` creates a todo from `{"title": "...", "done": false}` and returns it with its new ID.
- `GET /todos/{id}` returns one todo.
- `PUT /todos/{id}` replaces a todo's title/done flag.

It uses the Go 1.22 method+wildcard patterns on `http.ServeMux`, so no router dependency.

### How I tested it

Unit tests call the handler directly with `httptest.NewRecorder`, plus one test that starts a real
`httptest.NewServer` and uses `http.Get`.

Manually tried it with curl and it seemed fine. Some of the tests fail on my machine but I think the
tests are just too strict. Can you take a look before I ask for the real review?
