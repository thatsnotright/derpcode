# Review notes: order-service

## Planted issues

1. **`orders.go: Place` - money as `float64`.** `0.10 * 3 == 0.30000000000000004`, so totals and charges are off and `==` comparisons fail; rounding errors accumulate and discounts (`* 0.9`) are lossy. Fix: integer minor units (`PriceCents int64`), do the discount with integer math, format only at the edge (or use a decimal type). Caught by: `TestPlaceComputesTotal`.
2. **`Place` wraps with `%v`: `fmt.Errorf("reserve %s: %v", ...)`.** `%v` flattens the error to text, so `errors.Is(err, ErrOutOfStock)` fails even though the message looks right (the engineer's confusion in the PR). Fix: `%w`. Caught by: `TestPlaceOutOfStockReleasesEarlierItems`.
3. **`IsNotFound` compares with `==`** against a sentinel that `Get` wraps with `%w`. Fix: `errors.Is`. (Also look for `err == io.EOF`-style comparisons elsewhere and `err.Error() == "..."`.) Caught by: `TestGetNotFound`.
4. **`Place` swallows the payment error and returns bare `ErrPaymentDeclined`.** A gateway timeout is reported as "declined", the real cause is lost, and callers may tell the user their card was declined. Fix: wrap the cause: `fmt.Errorf("charge customer %s: %w", id, err)`; the gateway returns `ErrPaymentDeclined` for real declines. Caught by: `TestPlacePaymentFailureKeepsCause`.
5. **Context not propagated or checked.** `Charge` is called with `context.Background()` instead of `ctx`, and `Place` never checks `ctx.Err()`, so canceled requests still reserve stock and charge cards. Fix: pass `ctx` through, check `ctx.Err()` at entry. Rollback is the one place to use `context.WithoutCancel(ctx)` so cleanup still runs after cancellation. Caught by: `TestPlaceHonoursCanceledContext`.
6. **`Place` panics on an empty order.** `panic("orders: order has no items")` for bad user input takes down the whole server; return `ErrInvalidOrder` (wrapped with detail) instead. Validation is also incomplete (qty <= 0, negative price, empty SKU, too many lines). Caught by: `TestPlaceInvalidOrder` (the panic fails the test binary).
7. **Review only: ignored errors in rollback.** `s.inv.Release(...)` return values are dropped in `Place` and `Cancel`, so a failed release silently leaks reserved stock, and `Cancel` marks the order cancelled anyway. Fix: collect with `errors.Join` and return/log; do not mark cancelled if release failed.
8. **Review only: magic numbers and strings.** `100`, `0.9`, `"paid"`, `"cancelled"`, `"ord-%d"`. Fix: named constants (`bulkDiscountThresholdCents`, `bulkDiscountPercent`, `StatusPaid`), maybe a typed `Status`.
9. **Review only: not safe for concurrent use.** `orders` map and ID generation (`len(s.orders)+1`) are unguarded - a data race under `-race` with concurrent `Place`, and ID reuse if orders are ever deleted. `Get` also returns the internal `*Order` (callers can mutate it). Fix: mutex, counter, return a copy. Reference adds `TestConcurrentPlace`.

## Coaching / best-practice talking points

- Error wrapping: `%w` to keep the chain, `errors.Is`/`errors.As` to inspect, `errors.Join` for multiple failures; add context ("what were we doing") but do not repeat the underlying message.
- Sentinel vs typed errors; where to translate errors into API status codes.
- Panics are for programmer errors/impossible states, not bad input.
- Never store money in floats. Discuss rounding rules for percentage discounts (truncation here; be explicit).
- `context.Context` first parameter; propagate everywhere; check cancellation before side effects; `context.WithoutCancel` for cleanup.
- Compensation/rollback ordering, and idempotency (cancel twice; retry of Place) - a charge that succeeds but whose response is lost.
- Test design: fakes with behaviour you can toggle, table-driven validation tests, assert on `errors.Is` not strings, test cleanup paths and ordering.
- Exported identifiers need doc comments; named constants; `gofmt`/`go vet`.

## How to run the conversation

- Start with their confusion: "the message looks right but the test fails". Ask them to explain what `%v` vs `%w` stores in the error value.
- Ask "what happens if the client disconnects halfway through Place?" to surface ctx and rollback.
- Ask "what could the caller do differently for a timeout vs a decline?" to surface swallowed causes.
- Let them list the magic numbers themselves; time-box the money/float discussion to a few minutes and move on to design trade-offs (int cents vs decimal library).
- Concurrency and the API shape change (cents) are good "if we have time" topics.
