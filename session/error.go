package session

import (
	"net/http"

	"github.com/labstack/echo"
)

//Error -  common json error
type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	trace       string
}

//InternalServerError - internal server error
func InternalServerError(c echo.Context, e error) Error {

	errMsg := Error{
		Status:      http.StatusInternalServerError,
		Code:        http.StatusInternalServerError,
		Description: e.Error(),
	}

	return errMsg
}

//BadDataRequest - validator response
func BadDataRequest(c echo.Context, e error) Error {

	errMsg := Error{
		Status:      http.StatusBadRequest,
		Code:        http.StatusBadRequest,
		Description: e.Error(),
	}

	return errMsg
}

//InvalidLogin invalid login
func InvalidLogin(c echo.Context, e error) Error {

	errMsg := Error{
		Status:      http.StatusUnauthorized,
		Code:        http.StatusUnauthorized,
		Description: "Invalid Email / Password",
	}

	return errMsg
}
