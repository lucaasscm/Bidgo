package api

import (
	"errors"
	"net/http"

	"github.com/lucaasscm/bidgo/internal/jsonutils"
	"github.com/lucaasscm/bidgo/internal/services"
	"github.com/lucaasscm/bidgo/internal/usecase/bid"
)

func (api *Api) handleCreateBid(w http.ResponseWriter, r *http.Request) {
	productID, ok := parseProductID(w, r)
	if !ok {
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[bid.PlaceBidReq](r)
	if err != nil {
		if problems == nil {
			_ = jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid json body",
			})
			return
		}

		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	bidderID, ok := api.authenticatedUserID(w, r)
	if !ok {
		return
	}

	placedBid, err := api.BidsService.PlaceBid(r.Context(), productID, bidderID, data.BidAmount)
	if err != nil {
		api.writeBidError(w, r, err)
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, placedBid)
}

func (api *Api) handleListProductBids(w http.ResponseWriter, r *http.Request) {
	productID, ok := parseProductID(w, r)
	if !ok {
		return
	}

	bids, err := api.BidsService.ListBidsByProduct(r.Context(), productID)
	if err != nil {
		api.writeBidError(w, r, err)
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, bids)
}

// writeBidError maps bid service errors to their HTTP responses.
func (api *Api) writeBidError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, services.ErrProductNotFound):
		_ = jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "product not found",
		})
	case errors.Is(err, services.ErrCannotBidOwnProduct):
		_ = jsonutils.EncodeJson(w, r, http.StatusForbidden, map[string]any{
			"error": "you cannot bid on your own product",
		})
	case errors.Is(err, services.ErrAuctionEnded):
		_ = jsonutils.EncodeJson(w, r, http.StatusConflict, map[string]any{
			"error": "auction has already ended",
		})
	case errors.Is(err, services.ErrBidTooLow):
		_ = jsonutils.EncodeJson(w, r, http.StatusConflict, map[string]any{
			"error": "bid must be higher than the current highest bid and the base price",
		})
	default:
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
}
