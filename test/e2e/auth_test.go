package jwtauthv2_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jwtauthv2/api"
	"jwtauthv2/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var name = "kostya"
var password = "1234567"
var jsonUserString = fmt.Sprintf(`{ "name": "%v", "password": "%v" }`, name, password)

func TestSignUp(t *testing.T) {

	res := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBufferString(jsonUserString))
	require.Nil(t, err)

	Start.router.ServeHTTP(res, req)

	require.Equalf(t, http.StatusOK, res.Code, "Equal HTTP Status")

	user, err := Start.userStore.FindByName(name)
	require.Nil(t, err)

	responseData := new(api.JSendResponse)
	err = json.NewDecoder(res.Body).Decode(responseData)
	require.Nil(t, err)

	require.Equal(t, user.ID.String(), responseData.Data)
}

func TestSignIn(t *testing.T) {
	res := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/auth/signin", bytes.NewBufferString(jsonUserString))
	require.Nil(t, err)

	Start.router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	_, err = getTokenFromResBody(res)
	require.Nil(t, err)
}

func getTokenFromResBody(res *httptest.ResponseRecorder) (*entity.TokenPair, error) {
	responseBody := new(api.JSendResponse)
	tokenPair := new(entity.TokenPair)

	err := json.NewDecoder(res.Body).Decode(responseBody)

	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(&responseBody.Data)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(jsonData, tokenPair)

	return tokenPair, nil
}
