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
	ErrAuctionEnded        = errors.New("auction has already ended")
	ErrCannotBidOwnProduct = errors.New("seller cannot bid on their own product")
	ErrBidTooLow           = errors.New("bid amount is too low")
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (bs *BidsService) PlaceBid(ctx context.Context, productID, bidderID uuid.UUID, amount float64) (pgstore.Bid, error) {
	tx, err := bs.pool.Begin(ctx)
	if err != nil {
		return pgstore.Bid{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := bs.queries.WithTx(tx)

	// FOR UPDATE locks the product row so concurrent bids on the same product
	// are validated one at a time — otherwise two bids could both pass the
	// highest-bid check and both be accepted.
	product, err := qtx.GetProductByIdForUpdate(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, ErrProductNotFound
		}

		return pgstore.Bid{}, err
	}

	if product.SellerID == bidderID {
		return pgstore.Bid{}, ErrCannotBidOwnProduct
	}

	if product.IsSold || time.Now().After(product.AuctionEnd) {
		return pgstore.Bid{}, ErrAuctionEnded
	}

	highest, err := qtx.GetHighestBidByProductId(ctx, productID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		if amount < product.BasePrice {
			return pgstore.Bid{}, ErrBidTooLow
		}
	case err != nil:
		return pgstore.Bid{}, err
	default:
		if amount <= highest.BidAmount {
			return pgstore.Bid{}, ErrBidTooLow
		}
	}

	bid, err := qtx.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productID,
		BidderID:  bidderID,
		BidAmount: amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return pgstore.Bid{}, err
	}

	return bid, nil
}

func (bs *BidsService) ListBidsByProduct(ctx context.Context, productID uuid.UUID) ([]pgstore.Bid, error) {
	if _, err := bs.queries.GetProductById(ctx, productID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	bids, err := bs.queries.ListBidsByProductId(ctx, productID)
	if err != nil {
		return nil, err
	}

	if bids == nil {
		bids = []pgstore.Bid{}
	}

	return bids, nil
}
