package infrastructure

import (
	"context"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// DefaultGuestCartService defines the in memory guest cart service
	DefaultGuestCartService struct {
		defaultBehaviour *DefaultCartBehaviour
	}
)

var (
	_ cart.GuestCartService = (*DefaultGuestCartService)(nil)
)

// Inject dependencies
func (gcs *DefaultGuestCartService) Inject(
	InMemoryCartOrderBehaviour *DefaultCartBehaviour,
) {
	gcs.defaultBehaviour = InMemoryCartOrderBehaviour
}

// GetCart fetches a cart from the in memory guest cart service
func (gcs *DefaultGuestCartService) GetCart(ctx context.Context, cartID string) (*cart.Cart, error) {
	cart, err := gcs.defaultBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetNewCart gets a new cart from the in memory guest cart service
func (gcs *DefaultGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	guestCart := &cart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}

	err := gcs.defaultBehaviour.storeCart(guestCart)
	return guestCart, err
}

// GetModifyBehaviour returns the cart order behaviour of the service
func (gcs *DefaultGuestCartService) GetModifyBehaviour(context.Context) (cart.ModifyBehaviour, error) {
	return gcs.defaultBehaviour, nil
}

// RestoreCart restores a previously used guest cart
func (gcs *DefaultGuestCartService) RestoreCart(ctx context.Context, cart cart.Cart) (*cart.Cart, error) {
	guestCart := cart
	guestCart.ID = strconv.Itoa(rand.Int())

	err := gcs.defaultBehaviour.storeCart(&guestCart)
	return &guestCart, err
}
