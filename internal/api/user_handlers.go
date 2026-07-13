package api

import (
	"errors"
	"net/http"

	"github.com/lucaasscm/bidgo/internal/jsonutils"
	"github.com/lucaasscm/bidgo/internal/services"
	"github.com/lucaasscm/bidgo/internal/usecase/user"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
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

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}

		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleSignoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TO DO")
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	panic("TO DO")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TO DO")
}
