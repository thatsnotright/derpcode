## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add CSV user importer

### What it does

Adds `importer.Import` / `importer.ImportFiles`, which read user CSV exports and upsert them into an
in-memory `Store` keyed by user ID. Support will use it to load the nightly exports from the old
system.

Each file looks like:

```
id,name,email,age
1,Ada Lovelace,ada@example.com,36
2,Alan Turing,alan@example.com,41
```

- `Import(store, r)` imports one CSV stream.
- `ImportFiles(store, paths)` imports several files in order.
- Existing IDs are overwritten; new IDs are added.
- Bad input should be rejected and the support person should be told where the problem is, and a bad
  file should not leave the store half-updated.

### How I tested it

Tests feed small CSV strings to `Import` and check the store afterwards. The happy path works on all
the sample exports. I did not have a bad export to try, but the logic is simple. Can you take a look
before I ask for the real review?
