package routes

import (
	"net/http"

	"github.com/qaiswardag/go_backend_api_jwt/internal/controller/authcontroller"
	"github.com/qaiswardag/go_backend_api_jwt/internal/controller/homecontroller"
	"github.com/qaiswardag/go_backend_api_jwt/internal/controller/userregistercontroller"
	"github.com/qaiswardag/go_backend_api_jwt/internal/controller/usersessionscontroller"
	"github.com/qaiswardag/go_backend_api_jwt/internal/middleware"
)

type RouteHandler struct{}

func ChainMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func MainRouter() http.Handler {

	// TODO: Add methods for each route: "GET", "POST", "PUT" etc.
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		homecontroller.Show(w, r)
	}))

	mux.Handle("/user/sign-in", middleware.Cors(
		middleware.GlobalMiddleware(
			http.HandlerFunc(usersessionscontroller.Create),
		),
	))
	mux.Handle("/user/sign-up", middleware.Cors(
		middleware.GlobalMiddleware(
			http.HandlerFunc(userregistercontroller.Create),
		),
	))

	mux.Handle("/user/sign-out",
		middleware.Cors(
			middleware.GlobalMiddleware(
				middleware.RequireSessionMiddleware(
					http.HandlerFunc(authcontroller.Destroy),
				),
			),
		),
	)

	mux.Handle("/user/user", ChainMiddlewares(
		http.HandlerFunc(authcontroller.Show),
		middleware.RequireSessionMiddleware,
		middleware.GlobalMiddleware,
		middleware.Cors,
	))

	return mux
}
