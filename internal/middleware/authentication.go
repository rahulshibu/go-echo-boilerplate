package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.com/accubits/mapclub_multitenancy_admin/internal/configs"
)

var whitelist = [][2]string{
	{"GET", "^/api/v1/_hc"},
	{"POST", "^/api/v1/auth/login"},
	{"POST", "^/api/v1/auth/refresh"},
}

//Authorize - jwt validation
func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		//check if url is whitelisted
		if handleUnauthorized(c) {
			return next(c)
		}

		//get the token
		userToken := c.Request().Header.Get("Authorization")
		l := len("Bearer")
		if len(userToken) > l+1 && userToken[:l] == "Bearer" {
			userToken = userToken[l+1:]
		}

		//validate the token
		token, _ := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(configs.AppConfig.Secret), nil
		})

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"code":  http.StatusUnauthorized,
				"error": token.Claims.Valid().Error(),
			})
		}

		//else take the climes and keep in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			//claims
			name := claims["name"].(string)
			isAdmin := claims["admin"].(bool)
			userID := int(claims["id"].(float64))

			//log

			c.Set("name", name)
			c.Set("isAdmin", isAdmin)
			c.Set("userId", userID)

			return next(c)

		}

		return next(c)
	}
}

func handleUnauthorized(c echo.Context) bool {
	for _, pp := range whitelist {

		if pp[0] != c.Request().Method {
			continue
		}

		if matched, _ := regexp.MatchString(pp[1], strings.ToLower(c.Request().URL.Path)); matched {
			return true
		}
	}

	return false
}
