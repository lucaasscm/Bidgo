package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)
	api.Router.Use(api.Sessions.LoadAndSave)
	// dev-only (plain HTTP): mark requests as plaintext and allow the csrf
	// cookie without Secure — both must be removed/conditional behind HTTPS in production
	api.Router.Use(csrfPlaintextMiddleware)
	api.Router.Use(csrf.Protect(api.CSRFKey, csrf.Secure(false)))

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrftoken", api.handleGetCSRFToken)

			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/login", api.handleLoginUser)

				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/logout", api.handleLogoutUser)
				})
			})

			r.Route("/products", func(r chi.Router) {
				r.Get("/", api.handleListProducts)
				r.Get("/{id}", api.handleGetProduct)
				r.Get("/{id}/bids", api.handleListProductBids)

				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/", api.handleCreateProduct)
					r.Put("/{id}", api.handleUpdateProduct)
					r.Delete("/{id}", api.handleDeleteProduct)
					r.Post("/{id}/bids", api.handleCreateBid)
				})
			})
		})
	})
}

func csrfPlaintextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
	})
}
