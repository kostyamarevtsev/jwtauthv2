package api

import (
	"encoding/json"
	"fmt"
	"jwtauthv2"
	"jwtauthv2/entity"
	"jwtauthv2/usecase/auth"
	"net/http"

	jwt "github.com/form3tech-oss/jwt-go"
)

func HandleSignUp(authService *auth.AuthService) http.HandlerFunc {

	type requestBody struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(res http.ResponseWriter, req *http.Request) {
		data := new(requestBody)

		if err := json.NewDecoder(req.Body).Decode(data); err != nil {
			ResponseError(res, &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "Invalid Parse request body",
			})

			return
		}

		userid, err := authService.SignUp(data.Name, data.Password)

		if err != nil {
			ResponseError(res, err)
			return
		}

		Response(res, http.StatusOK, &JSendResponse{
			Status: APISUCCESS,
			Data:   userid.String(),
		})
	}
}

func HandleSignIn(authService *auth.AuthService) http.HandlerFunc {
	type requestBody struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(res http.ResponseWriter, req *http.Request) {
		data := new(requestBody)

		if err := json.NewDecoder(req.Body).Decode(data); err != nil {
			ResponseError(res, &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "Invalid Parse request body",
			})

			return
		}

		tokenPair, err := authService.SignIn(data.Name, data.Password)

		if err != nil {
			ResponseError(res, err)
			return
		}

		Response(res, http.StatusOK, &JSendResponse{
			Status: APISUCCESS,
			Data:   tokenPair,
		})
	}
}

func HandleRefresh(authService *auth.AuthService) http.HandlerFunc {
	type requestBody struct {
		RefreshToken string `json:"refreshToken"`
	}

	return func(res http.ResponseWriter, req *http.Request) {
		data := new(requestBody)

		if err := json.NewDecoder(req.Body).Decode(data); err != nil {
			ResponseError(res, &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "Invalid Parse request body",
			})

			return
		}

		tokenPair, err := authService.Refresh(data.RefreshToken)

		if err != nil {
			ResponseError(res, err)
			return
		}
		Response(res, http.StatusOK, &JSendResponse{
			Status: APISUCCESS,
			Data:   tokenPair,
		})
	}
}

func HandleLogout(authService *auth.AuthService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user := req.Context().Value("user")
		userid := user.(*jwt.Token).Claims.(jwt.MapClaims)["User"]

		id, err := entity.StringToID(fmt.Sprintf("%v", userid))

		if err != nil {
			ResponseError(res, &jwtauthv2.Error{
				Code:    jwtauthv2.EINVALID,
				Message: "Invalid Parse jwt claims",
			})
			return
		}

		err = authService.Logout(&id)

		if err != nil {
			ResponseError(res, err)
			return
		}

		Response(res, http.StatusOK, &JSendResponse{
			Status: APISUCCESS,
			Data:   nil,
		})

	}
}
