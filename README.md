# Interview Practice (local)

A CoderPad-style practice environment that runs entirely on your machine — no external
websites, no CDNs, no AI. It mirrors the two interview rounds:

| Round | What you practice | Problems |
|---|---|---|
| Software Development, part 1 | Implement the given methods in **Go** against tests | 7 (`problems/go-implement`) |
| Software Development, part 2 | Pair-review a junior engineer's **Go PR**: find bugs, make it pass `vet` + `-race`, coach on best practices | 5 (`problems/go-debug`) |
| Building Frontend Applications | Build an item-organizing **React** app step by step, with an accessibility checklist | 5 (`problems/react`) |

## Run it

Requirements: Node ≥ 23.6 (runs the TypeScript server directly) and Go ≥ 1.22.

```sh
npm start            # http://127.0.0.1:5173
```

The Go rounds need nothing else. The React live preview needs three local packages:

```sh
npm install          # esbuild, react, react-dom (one-time; needs network)
```

## Using it

- **Run** (`⌘/Ctrl+Enter`): Go problems run `gofmt -l`, `go vet` and `go test -race` and show per-test
  results; React problems re-bundle and reload the preview (it also refreshes as you type).
- **Format** (`⌘/Ctrl+S`): `gofmt` the current Go file.
- **Go docs (offline)**: `go doc strings.Builder` etc. — the interview allows looking up syntax.
- **Notes**: per-problem scratchpad — useful for review comments / coaching points in the PR round.
- **Reveal solution**: reference implementation, a diff against your code, and (for Go PRs) the
  reviewer notes listing every planted bug and the best-practice talking points.
- **Resume anytime**: your attempt at every problem autosaves to the browser's localStorage as you
  type (and on reload / tab close) — files, active tab, cursor/scroll position, last Go test run,
  notes and README checklist ticks. The last problem you had open reopens on reload, and problems
  you've started show a ● in the list. **Reset** restores the starter files for that problem.

## Adding problems

Each problem is a folder `problems/<go-implement|go-debug|react>/<id>/`:

- `meta.json` — `{"title", "difficulty", "minutes", "language": "go"|"react", "readonly": [files], "order"}`
- `README.md` — the statement shown in the sidebar
- starter files — everything else at the top level (Go: one package, `go.mod` is generated;
  React: `App.tsx` default-exporting the component, plus any `.tsx/.ts/.css`)
- `reference/` — solution files (same name replaces the starter file; new names are added)
- `REVIEW.md` — optional, hidden until reveal

`npm run verify` checks the bank: Go starters must build and fail, references must pass
`gofmt`/`vet`/`-race` (Go PRs 3× to catch flakes), and React starters/references must bundle.

## Layout

```
server/   Node http server (no dependencies): problem loader, Go runner, esbuild bundler
web/      Vanilla JS UI (no build step): editor, markdown, panels
problems/ The problem bank
scripts/  verify.ts
```
