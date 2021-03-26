package api

import (
	"jwtauthv2/usecase/auth"

	"github.com/gorilla/mux"
)

type Handler struct {
	authService *auth.AuthService
}

func NewHandler(authService *auth.AuthService) *Handler {
	return &Handler{authService}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", HandleSignUp(h.authService)).Methods("POST")
	auth.HandleFunc("/signin", HandleSignIn(h.authService)).Methods("POST")
	auth.HandleFunc("/refresh", HandleRefresh(h.authService)).Methods("POST")
	auth.Handle("/logout", JwtMiddleware(HandleLogout(h.authService))).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(JwtMiddleware)
	api.HandleFunc("/", HandleApi(h.authService)).Methods("GET")

	return router
}
