package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"

	"gitlab.com/accubits/mapclub_multitenancy_admin/internal/models"
	"gitlab.com/accubits/mapclub_multitenancy_admin/session"
)

type authHandler struct{}

func registerAuthHandlers(r *echo.Group) {

	h := authHandler{}

	r.POST("/auth/login", h.login)
	r.POST("/auth/refresh", h.token)

}

//authenticate user and generate token, refresh token
func (h *authHandler) login(c echo.Context) error {

	//auth model
	var authModel models.Auth

	var u models.User

	if err := c.Bind(&u); err != nil {
		return err
	}

	tokenResponse, err := authModel.Authenticate(u)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, session.InvalidLogin(c, err))
	}

	return c.JSON(http.StatusOK, tokenResponse)
}

//refresh token
func (h *authHandler) token(c echo.Context) error {
	//auth modek
	var authModel models.Auth

	var tr models.TokenReqBody

	if err := c.Bind(&tr); err != nil {
		return err
	}

	//refresh token
	tokenResponse, err := authModel.TokenRefresh(tr)

	if err != nil {
		log.Error(err.Error())
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, tokenResponse)

}
