package orders

import (
	"context"
	"errors"
	"testing"
)

var errGatewayTimeout = errors.New("gateway timeout")

type fakeInventory struct {
	stock    map[string]int
	reserved map[string]int
}

func newFakeInventory(stock map[string]int) *fakeInventory {
	return &fakeInventory{stock: stock, reserved: map[string]int{}}
}

func (f *fakeInventory) Reserve(_ context.Context, sku string, qty int) error {
	if f.stock[sku]-f.reserved[sku] < qty {
		return ErrOutOfStock
	}
	f.reserved[sku] += qty
	return nil
}

func (f *fakeInventory) Release(_ context.Context, sku string, qty int) error {
	f.reserved[sku] -= qty
	return nil
}

type fakePayments struct {
	err    error
	amount float64
	calls  int
}

func (f *fakePayments) Charge(ctx context.Context, _ string, amount float64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.calls++
	f.amount = amount
	return f.err
}

func newService(stock map[string]int, payErr error) (*Service, *fakeInventory, *fakePayments) {
	inv := newFakeInventory(stock)
	pay := &fakePayments{err: payErr}
	return NewService(inv, pay), inv, pay
}

func TestPlaceComputesTotal(t *testing.T) {
	svc, _, pay := newService(map[string]int{"pen": 10}, nil)
	o, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 3, Price: 0.10}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if o.Total != 0.30 {
		t.Fatalf("Total = %v; want 0.30", o.Total)
	}
	if pay.amount != 0.30 {
		t.Fatalf("charged %v; want 0.30", pay.amount)
	}
	if o.Status != "paid" {
		t.Fatalf("Status = %q; want paid", o.Status)
	}
}

func TestPlaceBulkDiscount(t *testing.T) {
	svc, _, _ := newService(map[string]int{"desk": 5}, nil)
	o, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "desk", Qty: 2, Price: 80}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if o.Total != 144 {
		t.Fatalf("Total = %v; want 144", o.Total)
	}
}

func TestPlaceOutOfStockReleasesEarlierItems(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10, "ink": 0}, nil)
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 2, Price: 1}, {SKU: "ink", Qty: 1, Price: 5}},
	})
	if !errors.Is(err, ErrOutOfStock) {
		t.Fatalf("err = %v; want ErrOutOfStock", err)
	}
	if inv.reserved["pen"] != 0 {
		t.Fatalf("pen still reserved: %d", inv.reserved["pen"])
	}
}

func TestPlacePaymentFailureKeepsCause(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10}, errGatewayTimeout)
	_, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 2, Price: 1}},
	})
	if !errors.Is(err, errGatewayTimeout) {
		t.Fatalf("err = %v; want it to wrap the gateway error", err)
	}
	if inv.reserved["pen"] != 0 {
		t.Fatalf("stock not released after failed payment: %d", inv.reserved["pen"])
	}
}

func TestPlaceInvalidOrder(t *testing.T) {
	svc, _, _ := newService(map[string]int{"pen": 10}, nil)
	_, err := svc.Place(context.Background(), Order{CustomerID: "c1"})
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("err = %v; want ErrInvalidOrder", err)
	}
}

func TestPlaceHonoursCanceledContext(t *testing.T) {
	svc, inv, pay := newService(map[string]int{"pen": 10}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Place(ctx, Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 1, Price: 1}},
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v; want context.Canceled", err)
	}
	if pay.calls != 0 || inv.reserved["pen"] != 0 {
		t.Fatalf("work was done despite canceled context: charges=%d reserved=%d", pay.calls, inv.reserved["pen"])
	}
}

func TestGetNotFound(t *testing.T) {
	svc, _, _ := newService(nil, nil)
	_, err := svc.Get(context.Background(), "ord-404")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false; want true", err)
	}
}

func TestCancel(t *testing.T) {
	svc, inv, _ := newService(map[string]int{"pen": 10}, nil)
	o, err := svc.Place(context.Background(), Order{
		CustomerID: "c1",
		Items:      []Item{{SKU: "pen", Qty: 4, Price: 1}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Cancel(context.Background(), o.ID); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.Get(context.Background(), o.ID)
	if got.Status != "cancelled" || inv.reserved["pen"] != 0 {
		t.Fatalf("status=%q reserved=%d", got.Status, inv.reserved["pen"])
	}
}
