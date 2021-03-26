package api_test

// Что возвращается 401 ошибка
// ЧТо при валидном токене вызывается handler и наоборот
// context value

import (
	"jwtauthv2/api"
	"jwtauthv2/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var user = entity.NewID()

func TestJWTMiddleware(t *testing.T) {
	user := entity.NewID()

	validToken, err := entity.IssueToken(&user, entity.AccessTokenTTL)
	require.Nil(t, err)

	expiredToken, err := entity.IssueToken(&user, 0)
	require.Nil(t, err)

	testTable := []struct {
		name, token        string
		expectedHTTPStatus int
	}{
		{"Middleware Expired Token", expiredToken, http.StatusUnauthorized},
		{"Middleware Valid Token", validToken, http.StatusOK},
	}

	for _, c := range testTable {
		t.Run(c.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/", nil)
			require.Nil(t, err)

			req.Header.Set("Authorization", "Bearer "+c.token)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.StatusOK, c.expectedHTTPStatus)
			})

			api.JwtMiddleware(handler).ServeHTTP(rr, req)
			require.Equal(t, rr.Code, c.expectedHTTPStatus)
		})
	}
}
