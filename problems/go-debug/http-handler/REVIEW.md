# Review notes: http-handler

## Planted issues

1. **`server.go: create` has no `return` after `http.Error` on bad JSON.** `http.Error` writes the response but does not stop the handler, so an empty todo is stored and a second body/`WriteHeader` is written ("superfluous response.WriteHeader" in logs). Fix: `return`. Caught by: `TestCreateBadJSON`.
2. **`server.go: create` responds 200 instead of 201 Created.** Fix: `http.StatusCreated`. Also the ReadAll-then-Unmarshal dance is unnecessary; decode straight from the (size-limited) body. Caught by: `TestCreate`.
3. **No `Content-Type: application/json` on any success response** (net/http sniffs and sends `text/plain`). Clients and proxies treat it wrongly. Fix: set the header before `WriteHeader`; factor a `writeJSON` helper. Caught by: `TestCreate`.
4. **`server.go: get` returns 400 for a missing todo.** Should be 404 Not Found (400 is for malformed requests). Caught by: `TestGetNotFound`.
5. **`server.go: decodeTodo` shadows `err`.** The named result `err` is shadowed by `if err := ...Decode(...)`, which only logs; the function always returns `nil` so `update` overwrites the todo with a zero value on bad JSON. Fix: do not shadow; return the decode error. Caught by: `TestUpdateBadJSON`. Good moment to mention `go vet -vettool` shadow analyzer / linters, and to avoid named results unless they help documentation.
6. **Review only: ignored errors and missing validation.** `json.NewEncoder(w).Encode(t)` result is ignored (at least log it); `io.ReadAll(r.Body)` has no size limit (use `http.MaxBytesReader`); no validation (empty/huge title accepted, unknown fields accepted); error bodies are plain text while success bodies are JSON (be consistent, e.g. `{"error": "..."}`).
7. **Review only (in the tests): the response body from `http.Get`/`http.Post` is never closed** in `TestRealServer`, and the errors from `http.Get` and `io.ReadAll` are discarded (`resp, _ = http.Get(...)` would nil-panic on failure). Fix: check `err`, `defer resp.Body.Close()`. (Server-side the request body is closed by the server, so no `defer r.Body.Close()` is needed in handlers; good thing to check they understand that.)

## Coaching / best-practice talking points

- Every `http.Error`/`writeError` is followed by `return`. Ask them to grep for it.
- Status code semantics: 201 on create, 400 vs 404 vs 422, 405 handled by the Go 1.22 mux for free.
- Centralise response writing (`writeJSON`, `writeError`) so headers/status/encoding are consistent and the encode error is handled in one place.
- Keep handlers thin: validate -> call store -> write response. Consider a `Store` interface so handlers can be tested with a fake; here the `Server` doubles as the store.
- Shadowed variables with `:=` inside `if`; prefer small functions and avoid named results.
- Test through the public `Routes()` handler; table-driven tests; `t.Helper()`; check headers and JSON bodies by decoding, not substring-matching, when exactness matters.
- `httptest.NewRecorder` for unit tests; `httptest.NewServer` only when transport behaviour matters. Close bodies in the client.
- Tests that failed locally were "too strict": discuss why that is a smell (the test was right).
- `gofmt`, doc comments, constants instead of magic numbers (`maxBodyBytes`).

## How to run the conversation

- Ask the engineer to run the tests and read the first failure aloud; let them decide what is wrong (tests vs code) and justify it.
- Ask "what does the client see if the JSON is bad?" and have them curl it mentally (status, headers, body).
- Ask "what happens to the code after `http.Error`?" instead of saying "missing return".
- Let them pick the order of fixes. Reserve validation/MaxBytesReader and body-closing for the "what else would you flag in review?" part.
