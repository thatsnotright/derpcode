# Bank Ledger

Implement an in-memory ledger that manages accounts, moves money between them, and keeps a per-account transaction history.

Amounts are positive integers in the smallest currency unit (cents), stored as `int64`. Balances never go negative.

## Methods to implement

```go
var (
	ErrAccountExists     = errors.New("account already exists")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
)

type TxType string

const (
	TxDeposit     TxType = "deposit"
	TxWithdraw    TxType = "withdraw"
	TxTransferOut TxType = "transfer_out"
	TxTransferIn  TxType = "transfer_in"
)

// Transaction is one entry in an account's history.
type Transaction struct {
	Type         TxType
	Amount       int64
	Counterparty string // other account for transfers, "" otherwise
}

func NewLedger() *Ledger

// CreateAccount opens a new account with a zero balance.
func (l *Ledger) CreateAccount(id string) error

// Deposit adds amount to the account.
func (l *Ledger) Deposit(id string, amount int64) error

// Withdraw removes amount from the account.
func (l *Ledger) Withdraw(id string, amount int64) error

// Transfer atomically moves amount from one account to another.
func (l *Ledger) Transfer(from, to string, amount int64) error

// Balance returns the current balance of the account.
func (l *Ledger) Balance(id string) (int64, error)

// History returns the account's transactions oldest first.
func (l *Ledger) History(id string) ([]Transaction, error)
```

## Rules

- Validate `amount` first (`ErrInvalidAmount` for `<= 0`), then accounts (`ErrAccountNotFound`, `ErrSameAccount`), then funds (`ErrInsufficientFunds`).
- A failed operation must change nothing: no balance change and no history entry.
- A transfer records `TxTransferOut` (Counterparty = destination) on the source and `TxTransferIn` (Counterparty = source) on the destination.
- `History` returns a copy; mutating it must not affect the ledger.

## Example

```go
l := NewLedger()
l.CreateAccount("alice")
l.CreateAccount("bob")
l.Deposit("alice", 1000)
l.Transfer("alice", "bob", 300)
l.Balance("alice")   // 700, nil
l.Balance("bob")     // 300, nil
l.Withdraw("bob", 500) // ErrInsufficientFunds
```

## Follow-ups

1. Make it safe for concurrent use. How do you avoid deadlock when two goroutines transfer in opposite directions?
2. Add idempotency keys so a retried transfer is applied only once.
3. Why is `float64` a bad choice for money? How would you support multiple currencies?
