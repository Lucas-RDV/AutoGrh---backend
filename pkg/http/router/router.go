package router

import (
	"net/http"

	"AutoGRH/pkg/Controller"
	middleware "AutoGRH/pkg/Controller/middleware"
	"AutoGRH/pkg/service"
)

func New(auth *service.AuthService) *http.ServeMux {
	mux := http.NewServeMux()
	authCtl := Controller.NewAuthController(auth)
	mux.HandleFunc("/auth/login", authCtl.Login)
	mux.Handle("/me", middleware.RequireAuth(auth, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if claims, ok := middleware.GetClaims(r.Context()); ok {
			w.Write([]byte("ola " + claims.Nome))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})))
	return mux
}
