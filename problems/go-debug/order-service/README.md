## Your task

You are reviewing a teammate's PR. Pair with them to get it PR-worthy:

- Find and fix the bugs in the code.
- Make the tests pass reliably with `go vet ./... && go test -race -count=1 ./...`.
- Note the best-practice feedback you would give in review (naming, docs, structure, tests).

---

## PR: Add order placement service

### What it does

Adds `orders.Service` with the checkout flow:

- `Place(ctx, order)` validates the order, applies the bulk discount (10% off orders over $100),
  reserves stock in the inventory service, charges the customer, and stores the paid order.
  If anything fails after stock was reserved, the reservations are released.
- `Get(ctx, id)` loads an order; `IsNotFound(err)` lets callers check for a missing order.
- `Cancel(ctx, id)` releases the stock and marks the order cancelled.

Inventory and payments are interfaces so they can be faked in tests. I defined sentinel errors
(`ErrOutOfStock`, `ErrNotFound`, ...) and wrap errors with context so logs are easier to read.

### How I tested it

Fake inventory and payment gateway in the test file; tests for totals, out-of-stock, a payment
failure, an invalid order, cancelling, and a canceled context.

Not all of them pass right now. The out-of-stock one fails even though the error message looks right
when I print it, which I do not understand.
