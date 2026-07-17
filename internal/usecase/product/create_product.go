package product

import (
	"context"
	"time"

	"github.com/lucaasscm/bidgo/internal/validator"
)

const minAuctionDuration = 2 * time.Hour

type CreateProductReq struct {
	ProductName string    `json:"product_name"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.ProductName, 255), "product_name", "this field cannot have more than 255 characters")
	eval.CheckField(req.BasePrice > 0, "base_price", "must be greater than zero")
	eval.CheckField(!req.AuctionEnd.IsZero(), "auction_end", "this field cannot be empty")
	eval.CheckField(!req.AuctionEnd.Before(time.Now().Add(minAuctionDuration)), "auction_end", "auction must end at least 2 hours from now")

	return eval
}
