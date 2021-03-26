package api_test

import (
	"encoding/json"
	"errors"
	"jwtauthv2"
	"jwtauthv2/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseError(t *testing.T) {

	testTable := []struct {
		name                   string
		inputErr               error
		expectedHTTPStatusCode int
		expectedBody           *api.JSendResponse
	}{
		{
			name: "EINVALID",
			inputErr: &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "EINVALID msg",
			},
			expectedHTTPStatusCode: http.StatusBadRequest,
			expectedBody: &api.JSendResponse{
				Status:  api.APIFAIL,
				Message: "EINVALID msg",
			},
		},
		{
			name: "EINTERNAL",
			inputErr: &jwtauthv2.Error{
				Op:  "internal op",
				Err: errors.New("test err"),
			},
			expectedHTTPStatusCode: http.StatusInternalServerError,
			expectedBody: &api.JSendResponse{
				Status:  api.APIERROR,
				Message: jwtauthv2.InternalMsg,
			},
		},
		{
			name: "ECONFLICT",
			inputErr: &jwtauthv2.Error{
				Code:    jwtauthv2.ECONFLICT,
				Message: "ECONFLICT msg",
			},
			expectedHTTPStatusCode: http.StatusConflict,
			expectedBody: &api.JSendResponse{
				Status:  api.APIFAIL,
				Message: "ECONFLICT msg",
			},
		},

		{
			name: "ENOTFOUND",
			inputErr: &jwtauthv2.Error{
				Code:    jwtauthv2.ENOTFOUND,
				Message: "ENOTFOUND msg",
			},
			expectedHTTPStatusCode: http.StatusNotFound,
			expectedBody: &api.JSendResponse{
				Status:  api.APIFAIL,
				Message: "ENOTFOUND msg",
			},
		},
		{
			name: "EPERMISSION",
			inputErr: &jwtauthv2.Error{
				Code:    jwtauthv2.EPERMISSION,
				Message: "EPERMISSION msg",
			},
			expectedHTTPStatusCode: http.StatusUnauthorized,
			expectedBody: &api.JSendResponse{
				Status:  api.APIFAIL,
				Message: "EPERMISSION msg",
			},
		},
	}

	for _, c := range testTable {
		t.Run(c.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			api.ResponseError(res, c.inputErr)
			require.Equalf(t, c.expectedHTTPStatusCode, res.Code, "expectedHTTPStatusCode")

			body := new(api.JSendResponse)
			err := json.NewDecoder(res.Body).Decode(body)
			require.Nil(t, err)

			require.Equalf(t, c.expectedBody.Status, body.Status, "expectedBody.Status")
			require.Equalf(t, c.expectedBody.Message, body.Message, "expectedBody.Message")
		})
	}

}
