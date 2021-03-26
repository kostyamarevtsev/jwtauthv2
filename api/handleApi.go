package api

import (
	"jwtauthv2/usecase/auth"
	"net/http"
)

func HandleApi(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
