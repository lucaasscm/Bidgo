package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucaasscm/bidgo/internal/store/pgstore"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrNotProductOwner = errors.New("user does not own this product")
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, sellerID uuid.UUID, productName string, basePrice float64, auctionEnd time.Time) (uuid.UUID, error) {
	args := pgstore.CreateProductParams{
		ProductName: productName,
		BasePrice:   basePrice,
		SellerID:    sellerID,
		AuctionEnd:  auctionEnd,
	}

	id, err := ps.queries.CreateProduct(ctx, args)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (ps *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (pgstore.Product, error) {
	product, err := ps.queries.GetProductById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Product{}, ErrProductNotFound
		}

		return pgstore.Product{}, err
	}

	return product, nil
}

func (ps *ProductService) ListProducts(ctx context.Context) ([]pgstore.Product, error) {
	products, err := ps.queries.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []pgstore.Product{}
	}

	return products, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, id, sellerID uuid.UUID, productName string, basePrice float64, auctionEnd time.Time, isSold bool) (pgstore.Product, error) {
	existing, err := ps.GetProductById(ctx, id)
	if err != nil {
		return pgstore.Product{}, err
	}

	if existing.SellerID != sellerID {
		return pgstore.Product{}, ErrNotProductOwner
	}

	args := pgstore.UpdateProductParams{
		ID:          id,
		ProductName: productName,
		BasePrice:   basePrice,
		AuctionEnd:  auctionEnd,
		IsSold:      isSold,
	}

	product, err := ps.queries.UpdateProduct(ctx, args)
	if err != nil {
		return pgstore.Product{}, err
	}

	return product, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id, sellerID uuid.UUID) error {
	existing, err := ps.GetProductById(ctx, id)
	if err != nil {
		return err
	}

	if existing.SellerID != sellerID {
		return ErrNotProductOwner
	}

	return ps.queries.DeleteProduct(ctx, id)
}
