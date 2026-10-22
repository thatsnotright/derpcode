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

type account struct {
	balance int64
	history []Transaction
}

// Ledger manages accounts and their transaction histories.
type Ledger struct {
	accounts map[string]*account
}

// NewLedger returns an empty ledger.
func NewLedger() *Ledger {
	return &Ledger{accounts: make(map[string]*account)}
}

// CreateAccount opens a new account with a zero balance.
func (l *Ledger) CreateAccount(id string) error {
	if _, ok := l.accounts[id]; ok {
		return ErrAccountExists
	}
	l.accounts[id] = &account{}
	return nil
}

// Deposit adds amount to the account.
func (l *Ledger) Deposit(id string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a, ok := l.accounts[id]
	if !ok {
		return ErrAccountNotFound
	}
	a.balance += amount
	a.history = append(a.history, Transaction{Type: TxDeposit, Amount: amount})
	return nil
}

// Withdraw removes amount from the account.
func (l *Ledger) Withdraw(id string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a, ok := l.accounts[id]
	if !ok {
		return ErrAccountNotFound
	}
	if a.balance < amount {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	a.history = append(a.history, Transaction{Type: TxWithdraw, Amount: amount})
	return nil
}

// Transfer atomically moves amount from one account to another.
func (l *Ledger) Transfer(from, to string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	src, ok := l.accounts[from]
	if !ok {
		return ErrAccountNotFound
	}
	dst, ok := l.accounts[to]
	if !ok {
		return ErrAccountNotFound
	}
	if from == to {
		return ErrSameAccount
	}
	if src.balance < amount {
		return ErrInsufficientFunds
	}
	src.balance -= amount
	dst.balance += amount
	src.history = append(src.history, Transaction{TxTransferOut, amount, to})
	dst.history = append(dst.history, Transaction{TxTransferIn, amount, from})
	return nil
}

// Balance returns the current balance of the account.
func (l *Ledger) Balance(id string) (int64, error) {
	a, ok := l.accounts[id]
	if !ok {
		return 0, ErrAccountNotFound
	}
	return a.balance, nil
}

// History returns the account's transactions, oldest first.
func (l *Ledger) History(id string) ([]Transaction, error) {
	a, ok := l.accounts[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	out := make([]Transaction, len(a.history))
	copy(out, a.history)
	return out, nil
}
