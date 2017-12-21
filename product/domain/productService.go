package domain

import (
	"context"
	"fmt"

	"go.aoe.com/flamingo/core/search/domain"
)

type (
	// ProductService interface
	ProductService interface {
		// Get a product
		Get(ctx context.Context, marketplaceCode string) (BasicProduct, error)
	}

	SearchResult struct {
		domain.Result
		Hits []BasicProduct
	}

	// SearchService is a typed search for products
	SearchService interface {
		Search(ctx context.Context, filter ...domain.Filter) (SearchResult, error)
	}

	// ProductNotFound is an error
	ProductNotFound struct {
		MarketplaceCode string
	}
)

// Error implements the error interface
func (err ProductNotFound) Error() string {
	return fmt.Sprintf("Product with Marketplace Code %q Not Found", err.MarketplaceCode)
}