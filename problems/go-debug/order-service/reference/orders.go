// Package orders implements order placement and cancellation.
package orders

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
)

// Sentinel errors returned (possibly wrapped) by the service.
var (
	ErrOutOfStock      = errors.New("out of stock")
	ErrNotFound        = errors.New("order not found")
	ErrPaymentDeclined = errors.New("payment declined")
	ErrInvalidOrder    = errors.New("invalid order")
)

// Order statuses.
const (
	StatusPaid      = "paid"
	StatusCancelled = "cancelled"
)

const (
	// Orders above this subtotal (in cents) get the bulk discount.
	bulkDiscountThresholdCents = 10_000
	// bulkDiscountPercent is taken off the subtotal of qualifying orders.
	bulkDiscountPercent = 10
	// maxLineItems bounds the number of lines in a single order.
	maxLineItems = 50
)

// Item is one line of an order. PriceCents is the unit price in cents.
type Item struct {
	SKU        string
	Qty        int
	PriceCents int64
}

// Order is a customer order. TotalCents is the amount charged, in cents.
type Order struct {
	ID         string
	CustomerID string
	Items      []Item
	TotalCents int64
	Status     string
}

// Inventory reserves and releases stock.
type Inventory interface {
	Reserve(ctx context.Context, sku string, qty int) error
	Release(ctx context.Context, sku string, qty int) error
}

// Payments charges customers.
type Payments interface {
	Charge(ctx context.Context, customerID string, amountCents int64) error
}

// Service places and cancels orders. It is safe for concurrent use.
type Service struct {
	inv Inventory
	pay Payments

	mu     sync.Mutex
	orders map[string]*Order
	nextID int
}

// NewService returns a Service using the given dependencies.
func NewService(inv Inventory, pay Payments) *Service {
	return &Service{inv: inv, pay: pay, orders: map[string]*Order{}}
}

func validate(o Order) error {
	switch {
	case o.CustomerID == "":
		return fmt.Errorf("%w: customer id is required", ErrInvalidOrder)
	case len(o.Items) == 0:
		return fmt.Errorf("%w: order has no items", ErrInvalidOrder)
	case len(o.Items) > maxLineItems:
		return fmt.Errorf("%w: more than %d items", ErrInvalidOrder, maxLineItems)
	}
	for _, it := range o.Items {
		if it.SKU == "" || it.Qty < 1 || it.PriceCents < 0 {
			return fmt.Errorf("%w: bad item %+v", ErrInvalidOrder, it)
		}
	}
	return nil
}

func totalCents(items []Item) int64 {
	var total int64
	for _, it := range items {
		total += it.PriceCents * int64(it.Qty)
	}
	if total > bulkDiscountThresholdCents {
		total -= total * bulkDiscountPercent / 100
	}
	return total
}

// release gives back reserved stock. It uses a context that is not canceled
// so rollback still happens when the caller's context is done.
func (s *Service) release(ctx context.Context, items []Item) error {
	ctx = context.WithoutCancel(ctx)
	var errs []error
	for _, it := range items {
		if err := s.inv.Release(ctx, it.SKU, it.Qty); err != nil {
			errs = append(errs, fmt.Errorf("release %s: %w", it.SKU, err))
		}
	}
	return errors.Join(errs...)
}

// Place validates, reserves stock for, and charges an order.
func (s *Service) Place(ctx context.Context, o Order) (*Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validate(o); err != nil {
		return nil, err
	}
	o.Items = slices.Clone(o.Items)
	total := totalCents(o.Items)

	reserved := make([]Item, 0, len(o.Items))
	for _, it := range o.Items {
		if err := s.inv.Reserve(ctx, it.SKU, it.Qty); err != nil {
			return nil, errors.Join(
				fmt.Errorf("reserve %s: %w", it.SKU, err),
				s.release(ctx, reserved),
			)
		}
		reserved = append(reserved, it)
	}

	if err := s.pay.Charge(ctx, o.CustomerID, total); err != nil {
		return nil, errors.Join(
			fmt.Errorf("charge customer %s: %w", o.CustomerID, err),
			s.release(ctx, reserved),
		)
	}

	o.TotalCents = total
	o.Status = StatusPaid

	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	o.ID = fmt.Sprintf("ord-%d", s.nextID)
	s.orders[o.ID] = &o
	return &o, nil
}

// Get returns a copy of the order with the given id.
func (s *Service) Get(ctx context.Context, id string) (*Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("get order %q: %w", id, ErrNotFound)
	}
	cp := *o
	cp.Items = slices.Clone(o.Items)
	return &cp, nil
}

// IsNotFound reports whether err means the order does not exist.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// Cancel releases the stock held by an order and marks it cancelled.
func (s *Service) Cancel(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return fmt.Errorf("cancel order %q: %w", id, ErrNotFound)
	}
	if o.Status == StatusCancelled {
		return nil
	}
	if err := s.release(ctx, o.Items); err != nil {
		return fmt.Errorf("cancel order %q: %w", id, err)
	}
	o.Status = StatusCancelled
	return nil
}
