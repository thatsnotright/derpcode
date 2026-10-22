package orders

import (
	"context"
	"errors"
	"sync"
	"testing"
)

var errGatewayTimeout = errors.New("gateway timeout")

type fakeInventory struct {
	mu         sync.Mutex
	stock      map[string]int
	reserved   map[string]int
	releaseErr error
}

func newFakeInventory(stock map[string]int) *fakeInventory {
	return &fakeInventory{stock: stock, reserved: map[string]int{}}
}

func (f *fakeInventory) Reserve(_ context.Context, sku string, qty int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stock[sku]-f.reserved[sku] < qty {
		return ErrOutOfStock
	}
	f.reserved[sku] += qty
	return nil
}

func (f *fakeInventory) Release(_ context.Context, sku string, qty int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.releaseErr != nil {
		return f.releaseErr
	}
	f.reserved[sku] -= qty
	return nil
}

func (f *fakeInventory) reservedQty(sku string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.reserved[sku]
}

type fakePayments struct {
	mu          sync.Mutex
	err         error
	amountCents int64
	calls       int
}

func (f *fakePayments) Charge(ctx context.Context, _ string, amountCents int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.amountCents = amountCents
	return f.err
}

func newService(stock map[string]int, payErr error) (*Service, *fakeInventory, *fakePayments) {
	inv := newFakeInventory(stock)
	pay := &fakePayments{err: payErr}
	return NewService(inv, pay), inv, pay
}

func TestPlaceTotals(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
		want  int64
	}{
		{"cents are exact", []Item{{SKU: "pen", Qty: 3, PriceCents: 10}}, 30},
		{"at threshold no discount", []Item{{SKU: "pen", Qty: 1, PriceCents: 10_000}}, 10_000},
		{"above threshold discounted", []Item{{SKU: "pen", Qty: 2, PriceCents: 8_000}}, 14_400},
		{"multiple lines", []Item{{SKU: "pen", Qty: 2, PriceCents: 199}, {SKU: "ink", Qty: 1, PriceCents: 501}}, 899},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, pay := newService(map[string]int{"pen": 10, "ink": 10}, nil)
			o, err := svc.Place(context.Background(), Order{CustomerID: "c1", Items: tt.items})
			if err != nil {
				t.Fatal(err)
			}
			if o.TotalCents != tt.want || pay.amountCents != tt.want {
				t.Fatalf("total=%d charged=%d; want %d", o.TotalCents, pay.amountCents, tt.want)
			}
			if o.Status != StatusPaid {
				t.Fatalf("Status = %q; want %q", o.Status, StatusPaid)
			}
		})
	}
}

func TestPlaceOutOfStockReleasesEarlierItems(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10, "ink": 0}, nil)
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 2, PriceCents: 100}, {SKU: "ink", Qty: 1, PriceCents: 500}},
	})
	if !errors.Is(err, ErrOutOfStock) {
		t.Fatalf("err = %v; want ErrOutOfStock", err)
	}
	if got := inv.reservedQty("pen"); got != 0 {
		t.Fatalf("pen still reserved: %d", got)
	}
}

func TestPlacePaymentFailureKeepsCause(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10}, errGatewayTimeout)
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 2, PriceCents: 100}},
	})
	if !errors.Is(err, errGatewayTimeout) {
		t.Fatalf("err = %v; want it to wrap the gateway error", err)
	}
	if got := inv.reservedQty("pen"); got != 0 {
		t.Fatalf("stock not released after failed payment: %d", got)
	}
}

func TestPlaceDeclinedIsDistinguishable(t *testing.T) {
	svc, _, _ := newService(map[string]int{"pen": 10}, ErrPaymentDeclined)
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 1, PriceCents: 100}},
	})
	if !errors.Is(err, ErrPaymentDeclined) {
		t.Fatalf("err = %v; want ErrPaymentDeclined", err)
	}
}

func TestPlaceReportsRollbackFailure(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10}, errGatewayTimeout)
	errRelease := errors.New("release failed")
	inv.releaseErr = errRelease
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 1, PriceCents: 100}},
	})
	if !errors.Is(err, errGatewayTimeout) || !errors.Is(err, errRelease) {
		t.Fatalf("err = %v; want both the charge and the release errors", err)
	}
}

func TestPlaceInvalidOrder(t *testing.T) {
	tests := []struct {
		name  string
		order Order
	}{
		{"no items", Order{CustomerID: "c1"}},
		{"no customer", Order{Items: []Item{{SKU: "pen", Qty: 1, PriceCents: 1}}}},
		{"zero qty", Order{CustomerID: "c1", Items: []Item{{SKU: "pen", Qty: 0, PriceCents: 1}}}},
		{"negative price", Order{CustomerID: "c1", Items: []Item{{SKU: "pen", Qty: 1, PriceCents: -1}}}},
		{"empty sku", Order{CustomerID: "c1", Items: []Item{{Qty: 1, PriceCents: 1}}}},
		{"too many lines", Order{CustomerID: "c1", Items: make([]Item, maxLineItems+1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, pay := newService(map[string]int{"pen": 10}, nil)
			_, err := svc.Place(context.Background(), tt.order)
			if !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("err = %v; want ErrInvalidOrder", err)
			}
			if pay.calls != 0 {
				t.Fatal("payment attempted for an invalid order")
			}
		})
	}
}

func TestPlaceHonoursCanceledContext(t *testing.T) {
	svc, inv, pay := newService(map[string]int{"pen": 10}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Place(ctx, Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 1, PriceCents: 100}},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v; want context.Canceled", err)
	}
	if pay.calls != 0 || inv.reservedQty("pen") != 0 {
		t.Fatalf("work was done despite canceled context: charges=%d reserved=%d", pay.calls, inv.reservedQty("pen"))
	}
}

func TestGetNotFound(t *testing.T) {
	svc, _, _ := newService(nil, nil)
	_, err := svc.Get(context.Background(), "ord-404")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false; want true", err)
	}
	if IsNotFound(errors.New("other")) {
		t.Fatal("IsNotFound(other) = true; want false")
	}
}

func TestCancel(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10}, nil)
	o, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 4, PriceCents: 100}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Cancel(context.Background(), o.ID); err != nil {
		t.Fatal(err)
	}
	// Cancelling twice must not release stock twice.
	if err := svc.Cancel(context.Background(), o.ID); err != nil {
		t.Fatal(err)
	}
	got, err := svc.Get(context.Background(), o.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusCancelled || inv.reservedQty("pen") != 0 {
		t.Fatalf("status=%q reserved=%d", got.Status, inv.reservedQty("pen"))
	}
}

func TestCancelNotFound(t *testing.T) {
	svc, _, _ := newService(nil, nil)
	if err := svc.Cancel(context.Background(), "nope"); !IsNotFound(err) {
		t.Fatalf("err = %v; want not found", err)
	}
}

func TestConcurrentPlace(t *testing.T) {
	svc, _, _ := newService(map[string]int{"pen": 1000}, nil)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Place(context.Background(), Order{
				CustomerID: "c1",
				Items:      []Item{{SKU: "pen", Qty: 1, PriceCents: 100}},
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}
