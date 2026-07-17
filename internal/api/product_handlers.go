package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lucaasscm/bidgo/internal/jsonutils"
	"github.com/lucaasscm/bidgo/internal/services"
	"github.com/lucaasscm/bidgo/internal/usecase/product"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)
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

	sellerID, ok := api.authenticatedUserID(w, r)
	if !ok {
		return
	}

	id, err := api.ProductService.CreateProduct(r.Context(), sellerID, data.ProductName, data.BasePrice, data.AuctionEnd)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"product_id": id,
	})
}

func (api *Api) handleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := api.ProductService.ListProducts(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, products)
}

func (api *Api) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseProductID(w, r)
	if !ok {
		return
	}

	prod, err := api.ProductService.GetProductById(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			_ = jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"error": "product not found",
			})
			return
		}

		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, prod)
}

func (api *Api) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseProductID(w, r)
	if !ok {
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[product.UpdateProductReq](r)
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

	sellerID, ok := api.authenticatedUserID(w, r)
	if !ok {
		return
	}

	prod, err := api.ProductService.UpdateProduct(r.Context(), id, sellerID, data.ProductName, data.BasePrice, data.AuctionEnd, data.IsSold)
	if err != nil {
		api.writeProductError(w, r, err)
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, prod)
}

func (api *Api) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := parseProductID(w, r)
	if !ok {
		return
	}

	sellerID, ok := api.authenticatedUserID(w, r)
	if !ok {
		return
	}

	if err := api.ProductService.DeleteProduct(r.Context(), id, sellerID); err != nil {
		api.writeProductError(w, r, err)
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "product deleted successfully",
	})
}

// parseProductID reads the "id" URL parameter as a UUID, writing a 400 response
// and returning ok=false when it is missing or malformed.
func parseProductID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
		return uuid.UUID{}, false
	}

	return id, true
}

// authenticatedUserID pulls the logged-in user's id from the session, writing a
// 401 response and returning ok=false when it is absent (routes using this must
// sit behind AuthMiddleware).
func (api *Api) authenticatedUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, ok := api.Sessions.Get(r.Context(), "authenticatedUserId").(uuid.UUID)
	if !ok {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
			"error": "must be logged in",
		})
		return uuid.UUID{}, false
	}

	return id, true
}

// writeProductError maps product service errors to their HTTP responses.
func (api *Api) writeProductError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, services.ErrProductNotFound):
		_ = jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "product not found",
		})
	case errors.Is(err, services.ErrNotProductOwner):
		_ = jsonutils.EncodeJson(w, r, http.StatusForbidden, map[string]any{
			"error": "you do not own this product",
		})
	default:
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
}
