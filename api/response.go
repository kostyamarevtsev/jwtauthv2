package api

import (
	"encoding/json"
	"jwtauthv2"
	"log"
	"net/http"
)

const APISUCCESS = "SUCCESS"
const APIFAIL = "FAIL"
const APIERROR = "ERROR"

type JSendResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseError(res http.ResponseWriter, err error) {
	appErrCode := jwtauthv2.ErrorCode(err)
	appErrMessage := jwtauthv2.ErrorMessage(err)

	if appErrCode == jwtauthv2.EINTERNAL {
		// Логируем
		log.Println(err.Error())
		Response(res, http.StatusInternalServerError, &JSendResponse{
			Status:  APIERROR,
			Message: appErrMessage,
		})
		return
	}

	var httpStatusCode int

	if appErrCode == jwtauthv2.ECONFLICT {
		httpStatusCode = http.StatusConflict
	}

	if appErrCode == jwtauthv2.ENOTFOUND {
		httpStatusCode = http.StatusNotFound
	}

	if appErrCode == jwtauthv2.EINVALID {
		httpStatusCode = http.StatusBadRequest
	}

	if appErrCode == jwtauthv2.EPERMISSION {
		httpStatusCode = http.StatusUnauthorized
	}

	Response(res, httpStatusCode, &JSendResponse{
		Status:  APIFAIL,
		Message: appErrMessage,
	})
}

func Response(res http.ResponseWriter, httpStatusCode int, body *JSendResponse) {
	res.WriteHeader(httpStatusCode)
	res.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(res).Encode(*body)
	if err != nil {
		log.Fatal(err)
	}
}
