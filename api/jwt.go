package api

import (
	"context"
	"jwtauthv2"
	"jwtauthv2/entity"
	"net/http"
	"strings"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := strings.Split(req.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			ResponseError(res, &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "Parse Authorization Header",
			})
			return
		}

		tokenString := authHeader[1]
		claims, err := entity.ParseToken(tokenString)

		if err != nil {
			ResponseError(res, err)
			return
		}

		ctx := context.WithValue(req.Context(), "user", claims.User)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
