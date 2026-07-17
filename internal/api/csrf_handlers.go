package api

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/lucaasscm/bidgo/internal/jsonutils"
)

func (api *Api) handleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": csrf.Token(r),
	})
}
