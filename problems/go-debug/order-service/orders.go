// Package orders implements order placement and cancellation.
package orders

import (
	"context"
	"errors"
	"fmt"
)

// Sentinel errors returned (possibly wrapped) by the service.
var (
	ErrOutOfStock      = errors.New("out of stock")
	ErrNotFound        = errors.New("order not found")
	ErrPaymentDeclined = errors.New("payment declined")
	ErrInvalidOrder    = errors.New("invalid order")
)

// Item is one line of an order. Price is in dollars.
type Item struct {
	SKU   string
	Qty   int
	Price float64
}

// Order is a customer order.
type Order struct {
	ID         string
	CustomerID string
	Items      []Item
	Total      float64
	Status     string
}

// Inventory reserves and releases stock.
type Inventory interface {
	Reserve(ctx context.Context, sku string, qty int) error
	Release(ctx context.Context, sku string, qty int) error
}

// Payments charges customers.
type Payments interface {
	Charge(ctx context.Context, customerID string, amount float64) error
}

// Service places and cancels orders.
type Service struct {
	inv    Inventory
	pay    Payments
	orders map[string]*Order
}

// NewService returns a Service using the given dependencies.
func NewService(inv Inventory, pay Payments) *Service {
	return &Service{inv: inv, pay: pay, orders: map[string]*Order{}}
}

// Place validates, reserves stock for, and charges an order.
func (s *Service) Place(ctx context.Context, o Order) (*Order, error) {
	if len(o.Items) == 0 {
		panic("orders: order has no items")
	}
	if o.CustomerID == "" {
		return nil, ErrInvalidOrder
	}

	var total float64
	for _, it := range o.Items {
		total += it.Price * float64(it.Qty)
	}
	if total > 100 {
		total = total * 0.9
	}

	var reserved []Item
	for _, it := range o.Items {
		if err := s.inv.Reserve(ctx, it.SKU, it.Qty); err != nil {
			for _, r := range reserved {
				s.inv.Release(ctx, r.SKU, r.Qty)
			}
			return nil, fmt.Errorf("reserve %s: %v", it.SKU, err)
		}
		reserved = append(reserved, it)
	}

	if err := s.pay.Charge(context.Background(), o.CustomerID, total); err != nil {
		for _, r := range reserved {
			s.inv.Release(ctx, r.SKU, r.Qty)
		}
		return nil, ErrPaymentDeclined
	}

	o.Total = total
	o.Status = "paid"
	o.ID = fmt.Sprintf("ord-%d", len(s.orders)+1)
	s.orders[o.ID] = &o
	return &o, nil
}

// Get returns the order with the given id.
func (s *Service) Get(ctx context.Context, id string) (*Order, error) {
	o, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("get order %q: %w", id, ErrNotFound)
	}
	return o, nil
}

// IsNotFound reports whether err means the order does not exist.
func IsNotFound(err error) bool {
	return err == ErrNotFound
}

// Cancel releases the stock held by an order and marks it cancelled.
func (s *Service) Cancel(ctx context.Context, id string) error {
	o, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	for _, it := range o.Items {
		s.inv.Release(ctx, it.SKU, it.Qty)
	}
	o.Status = "cancelled"
	return nil
}
