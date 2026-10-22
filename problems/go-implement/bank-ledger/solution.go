package solution

import "errors"

var (
	ErrAccountExists     = errors.New("account already exists")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
)

// TxType identifies the kind of a history entry.
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

// Ledger manages accounts and their transaction histories.
type Ledger struct {
	// TODO: add fields
}

// NewLedger returns an empty ledger.
func NewLedger() *Ledger {
	// TODO: implement
	return &Ledger{}
}

// CreateAccount opens a new account with a zero balance.
func (l *Ledger) CreateAccount(id string) error {
	// TODO: implement
	return nil
}

// Deposit adds amount to the account.
func (l *Ledger) Deposit(id string, amount int64) error {
	// TODO: implement
	return nil
}

// Withdraw removes amount from the account.
func (l *Ledger) Withdraw(id string, amount int64) error {
	// TODO: implement
	return nil
}

// Transfer atomically moves amount from one account to another.
func (l *Ledger) Transfer(from, to string, amount int64) error {
	// TODO: implement
	return nil
}

// Balance returns the current balance of the account.
func (l *Ledger) Balance(id string) (int64, error) {
	// TODO: implement
	return 0, nil
}

// History returns the account's transactions, oldest first.
func (l *Ledger) History(id string) ([]Transaction, error) {
	// TODO: implement
	return nil, nil
}
