package solution

import (
	"errors"
	"reflect"
	"testing"
)

func newFunded(t *testing.T) *Ledger {
	t.Helper()
	l := NewLedger()
	for _, id := range []string{"alice", "bob"} {
		if err := l.CreateAccount(id); err != nil {
			t.Fatalf("CreateAccount(%q) = %v", id, err)
		}
	}
	if err := l.Deposit("alice", 1000); err != nil {
		t.Fatalf("Deposit = %v", err)
	}
	return l
}

func wantBalance(t *testing.T, l *Ledger, id string, want int64) {
	t.Helper()
	got, err := l.Balance(id)
	if err != nil || got != want {
		t.Errorf("Balance(%q) = (%d, %v), want (%d, nil)", id, got, err, want)
	}
}

func TestCreateAccount(t *testing.T) {
	l := NewLedger()
	if err := l.CreateAccount("a"); err != nil {
		t.Fatalf("first CreateAccount = %v", err)
	}
	if err := l.CreateAccount("a"); !errors.Is(err, ErrAccountExists) {
		t.Errorf("duplicate CreateAccount = %v, want ErrAccountExists", err)
	}
	wantBalance(t, l, "a", 0)
}

func TestBalanceUnknownAccount(t *testing.T) {
	l := NewLedger()
	if _, err := l.Balance("ghost"); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Balance = %v, want ErrAccountNotFound", err)
	}
	if _, err := l.History("ghost"); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("History = %v, want ErrAccountNotFound", err)
	}
}

func TestDepositAndWithdraw(t *testing.T) {
	l := newFunded(t)
	if err := l.Withdraw("alice", 250); err != nil {
		t.Fatalf("Withdraw = %v", err)
	}
	wantBalance(t, l, "alice", 750)
	if err := l.Withdraw("alice", 750); err != nil {
		t.Fatalf("Withdraw to zero = %v", err)
	}
	wantBalance(t, l, "alice", 0)
}

func TestInvalidAmounts(t *testing.T) {
	l := newFunded(t)
	for _, amt := range []int64{0, -5} {
		if err := l.Deposit("alice", amt); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Deposit(%d) = %v, want ErrInvalidAmount", amt, err)
		}
		if err := l.Withdraw("alice", amt); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Withdraw(%d) = %v, want ErrInvalidAmount", amt, err)
		}
		if err := l.Transfer("alice", "bob", amt); !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Transfer(%d) = %v, want ErrInvalidAmount", amt, err)
		}
	}
	wantBalance(t, l, "alice", 1000)
}

func TestUnknownAccountErrors(t *testing.T) {
	l := newFunded(t)
	if err := l.Deposit("ghost", 10); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Deposit = %v, want ErrAccountNotFound", err)
	}
	if err := l.Withdraw("ghost", 10); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Withdraw = %v, want ErrAccountNotFound", err)
	}
	if err := l.Transfer("ghost", "bob", 10); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Transfer from ghost = %v, want ErrAccountNotFound", err)
	}
	if err := l.Transfer("alice", "ghost", 10); !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Transfer to ghost = %v, want ErrAccountNotFound", err)
	}
	wantBalance(t, l, "alice", 1000)
}

func TestInsufficientFundsChangesNothing(t *testing.T) {
	l := newFunded(t)
	if err := l.Withdraw("alice", 1001); !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Withdraw = %v, want ErrInsufficientFunds", err)
	}
	if err := l.Transfer("alice", "bob", 1001); !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Transfer = %v, want ErrInsufficientFunds", err)
	}
	wantBalance(t, l, "alice", 1000)
	wantBalance(t, l, "bob", 0)
	h, _ := l.History("alice")
	if len(h) != 1 {
		t.Errorf("failed operations must not be recorded; history = %v", h)
	}
	if h, _ := l.History("bob"); len(h) != 0 {
		t.Errorf("bob history = %v, want empty", h)
	}
}

func TestTransfer(t *testing.T) {
	l := newFunded(t)
	if err := l.Transfer("alice", "bob", 300); err != nil {
		t.Fatalf("Transfer = %v", err)
	}
	wantBalance(t, l, "alice", 700)
	wantBalance(t, l, "bob", 300)
}

func TestTransferToSelf(t *testing.T) {
	l := newFunded(t)
	if err := l.Transfer("alice", "alice", 10); !errors.Is(err, ErrSameAccount) {
		t.Errorf("Transfer = %v, want ErrSameAccount", err)
	}
	wantBalance(t, l, "alice", 1000)
}

func TestHistory(t *testing.T) {
	l := newFunded(t)
	l.Withdraw("alice", 100)
	l.Transfer("alice", "bob", 300)
	l.Deposit("bob", 5)

	wantAlice := []Transaction{
		{TxDeposit, 1000, ""},
		{TxWithdraw, 100, ""},
		{TxTransferOut, 300, "bob"},
	}
	wantBob := []Transaction{
		{TxTransferIn, 300, "alice"},
		{TxDeposit, 5, ""},
	}
	if got, err := l.History("alice"); err != nil || !reflect.DeepEqual(got, wantAlice) {
		t.Errorf("History(alice) = %v, %v; want %v", got, err, wantAlice)
	}
	if got, err := l.History("bob"); err != nil || !reflect.DeepEqual(got, wantBob) {
		t.Errorf("History(bob) = %v, %v; want %v", got, err, wantBob)
	}
}

func TestHistoryReturnsCopy(t *testing.T) {
	l := newFunded(t)
	h, _ := l.History("alice")
	if len(h) != 1 {
		t.Fatalf("history length = %d, want 1", len(h))
	}
	h[0].Amount = 999999
	h2, _ := l.History("alice")
	if h2[0].Amount != 1000 {
		t.Errorf("mutating returned history changed ledger: %v", h2)
	}
}
