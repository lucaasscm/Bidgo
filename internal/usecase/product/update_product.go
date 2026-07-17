package product

import (
	"context"
	"time"

	"github.com/lucaasscm/bidgo/internal/validator"
)

type UpdateProductReq struct {
	ProductName string    `json:"product_name"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
	IsSold      bool      `json:"is_sold"`
}

func (req UpdateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.ProductName, 255), "product_name", "this field cannot have more than 255 characters")
	eval.CheckField(req.BasePrice > 0, "base_price", "must be greater than zero")
	eval.CheckField(!req.AuctionEnd.IsZero(), "auction_end", "this field cannot be empty")

	return eval
}
