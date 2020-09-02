package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

//health

type health struct {
	Build string `json:"build"`
}

//RegisterRoutes - golbal route register
func RegisterRoutes(e *echo.Echo) {

	//api group
	api := e.Group("/api/v1")
	api.GET("/_hc", pingHealth)

	//register routes
	registerAuthHandlers(api)

	registerUserHandler(api)

}

func pingHealth(c echo.Context) error {

	u := &health{
		Build: "1.0",
	}

	return c.JSON(http.StatusOK, u)
}

//RegisterErrorHanders - register error
func RegisterErrorHanders(e *echo.Echo) {

	//register not found
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.JSON(report.Code, report)
	}

}
