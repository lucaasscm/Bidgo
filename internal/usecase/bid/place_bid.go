package bid

import (
	"context"

	"github.com/lucaasscm/bidgo/internal/validator"
)

type PlaceBidReq struct {
	BidAmount float64 `json:"bid_amount"`
}

func (req PlaceBidReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(req.BidAmount > 0, "bid_amount", "must be greater than zero")

	return eval
}
