# Review notes: csv-importer

## Planted issues

1. **`importer.go: Import` skips the first line with a bare `sc.Scan() // skip header`.** The header is never validated: a file without a header silently loses its first user, a file with columns in a different order imports garbage (name as id, id as name), and an empty file is "successfully" imported. Fix: read the header, compare it to the expected columns (trim/case-fold if you want to be lenient), and return an error that says "header". Caught by: `TestImportRejectsBadHeader`.
2. **`importer.go: Import` never checks `sc.Err()` and uses the default `bufio.Scanner` buffer (64 KiB max token).** One long line makes `Scan` return false with `bufio.ErrTooLong`; because the error is never read, the import "succeeds" and every row after that line is silently dropped. Fix: use `encoding/csv` (reads arbitrarily long records), or at least `sc.Buffer(...)` and `if err := sc.Err(); err != nil`. Caught by: `TestImportHandlesVeryLongLines`.
3. **`importer.go: Import` ignores both `strconv.Atoi` errors (`id, _ := ...`).** A bad id or age becomes `0`: `abc` overwrites user 0, `old` makes someone age 0. Fix: check and return the errors, wrapped with `%w` and the field name. Caught by: `TestImportReportsLineNumbersAndWrapsCause`.
4. **`importer.go: Import` returns `errors.New("invalid row")`, and returns the bare `err` from `openFile`.** No line number, no file name, no cause, so support cannot find the bad row in a 50k-line export and callers cannot `errors.Is` anything. Fix: `fmt.Errorf("line %d: %w", n, err)` in `Import` and `fmt.Errorf("%s: %w", path, err)` in `ImportFiles`. Caught by: `TestImportReportsShortRows`, `TestImportReportsLineNumbersAndWrapsCause`, `TestImportFilesNamesTheBadFileAndAppliesNothing`.
5. **`importer.go: Import` upserts each row as it is parsed (and `ImportFiles` applies file by file).** When row N fails, rows 1..N-1 (and earlier files) are already in the store, and existing users have been overwritten: the state is half-updated and hard to recover from. Fix: parse and validate everything into a slice first, then apply with one `UpsertAll` under a single lock. Caught by: `TestImportIsAllOrNothing`, `TestImportFilesNamesTheBadFileAndAppliesNothing`.
6. **`importer.go: ImportFiles` uses `defer f.Close()` inside the loop.** Deferred calls run when the *function* returns, so with 5,000 daily exports all 5,000 files are open at once (hits `ulimit -n`: "too many open files"). Fix: move per-file work into a helper (`parseFile`) so the `defer` runs per file, or close explicitly. Caught by: `TestImportFilesClosesEachFileBeforeOpeningTheNext`.
7. **Review only: hand-rolled CSV parsing with `strings.Split(line, ",")`.** Quoted fields (`"Lovelace, Ada"`), escaped quotes and quoted newlines break; a single comma in a name turns a valid row into "invalid row". Use `encoding/csv` (the reference does; `TestImportHandlesQuotedFields` exists only in the reference). The scanner also hides the "which line am I on" bookkeeping that `csv.Reader.FieldPos` gives you for free.
8. **Review only: no semantic validation or duplicate handling.** Nothing checks that email looks like an email, that age is in a sane range, that the id is positive, or that the same id does not appear twice in one file (silent last-wins). Also `Import` takes a concrete `*Store`; an interface (`Upserter`) would make it easier to test and to swap for a database. Decide what is a hard error vs a warning and say so in the doc comment.

## Coaching / best-practice talking points

- Never write `x, _ := strconv...` for external input. Ask what the "default" value means to the business (user 0? age 0?).
- Every `Scanner` loop ends with `sc.Err()`. Know the 64 KiB limit; know `csv.Reader` / `bufio.Reader` as the alternatives for unbounded lines.
- `defer` is function-scoped. Loops that open resources need a helper function or explicit close. Same advice for locks, transactions and HTTP bodies.
- Validate first, then commit (parse, then apply). It is the cheap version of a transaction; for a real DB use a transaction. Make the store API support batches (`UpsertAll`) so atomicity is under one lock.
- Error design: add context at each layer (file, line, field), wrap with `%w` so callers can use `errors.Is`/`errors.As`, and do not repeat the same context twice. Mention collecting multiple errors (`errors.Join`) so users fix all rows in one pass, and when fail-fast is better.
- Testability: the `openFile` variable is a pragmatic seam; mention `io/fs.FS` as the cleaner one (`ImportFS(store, fsys, names...)` works with `os.DirFS` and `fstest.MapFS`). Remind them to restore globals with `t.Cleanup` and that such tests cannot use `t.Parallel`.
- Tests: table-driven with `t.Run`, assert on error text sparingly (line numbers are part of the contract here), and test the boundaries (empty input, header only, very long line).
- Possible extensions to discuss: streaming for huge files, progress reporting, per-row warnings, dry-run mode, concurrency across files (and then why the batch/lock design matters).

## How to run the conversation

- Ask the engineer to feed the importer a bad file by hand ("what does support see when row 5 is broken?") before looking at the tests. Let the failure output guide them.
- For the header: ask what `sc.Scan()` returns and what happens with an empty file or a headerless file; let them find all three cases.
- For the scanner: ask "how many users do I get from a file with 3 rows if row 2 is 100 KB?" and have them predict, then run it. Then ask how they would have noticed without the test.
- When they fix error handling, ask what the person reading the message will do next; they should arrive at line numbers and file names themselves.
- Ask "what is in the store if the import fails halfway?" and let them design the parse-then-apply split.
- The `defer` in the loop is a good review-by-reading catch; ask how many files are open at the end of the loop before revealing the test.
- Keep the `strings.Split` vs `encoding/csv` and validation/duplicate topics for the "what else would you flag?" part.
