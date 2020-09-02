package models

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/accubits/mapclub_multitenancy_admin/internal/configs"
	"gitlab.com/accubits/mapclub_multitenancy_admin/internal/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

//Auth - authentication handler
type Auth struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

//TokenReqBody - refresh token struct
type TokenReqBody struct {
	RefreshToken string `json:"refresh_token"`
}

//Authenticate user with user details
func (a *Auth) Authenticate(u User) (Auth, error) {
	var (
		user User
	)

	//initialize the jwt
	token := jwt.New(jwt.SigningMethodHS256)

	//check user login correct or not
	//TODO :
	dbconn := database.GetSharedConnection()

	//hash password
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	// u.PasswordHash = string(hashedPassword)

	err := dbconn.Debug().Where("email = ?", u.Email).First(&user).Error

	if err != nil {
		return Auth{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password))

	if err != nil {
		return Auth{}, err
	}

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Email // replace with username
	claims["admin"] = user.IsAdmin
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(configs.AppConfig.Secret))
	if err != nil {
		log.Error(err.Error())
		return Auth{}, err
	}

	// Generate encoded refrest token and send it as response.
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["name"] = user.Email // replace with username
	rtClaims["admin"] = user.IsAdmin
	rtClaims["sub"] = 1
	rtClaims["id"] = user.ID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString([]byte(configs.AppConfig.Secret))
	if err != nil {
		log.Error(err.Error())
		return Auth{}, err
	}

	auth := Auth{
		Token:        t,
		RefreshToken: rt,
	}

	return auth, nil
}

//TokenRefresh - refresh token
func (a *Auth) TokenRefresh(tr TokenReqBody) (Auth, error) {

	//validate the token and alg tag
	token, _ := jwt.Parse(tr.RefreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Unexpected signing method")
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.AppConfig.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := a.generateTokenPair(claims)
			if err != nil {
				log.Error(err.Error())
				return Auth{}, err
			}

			return newTokenPair, nil
		}

		return Auth{}, errors.New("Invalid Claim")
	}

	return Auth{}, nil
}

func (a *Auth) generateTokenPair(c jwt.MapClaims) (Auth, error) {
	//initialize the jwt
	token := jwt.New(jwt.SigningMethodHS256)

	//check user login correct or not
	//TODO :

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = c["name"] // replace with username
	claims["admin"] = c["admin"]
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID          works too)
	t, err := token.SignedString([]byte(configs.AppConfig.Secret))
	if err != nil {
		log.Error(err.Error())
		return Auth{}, err
	}

	// Generate encoded refrest token and send it as response.
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["name"] = c["name"] // replace with username
	rtClaims["admin"] = c["admin"]
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString([]byte(configs.AppConfig.Secret))
	if err != nil {
		log.Error(err.Error())
		return Auth{}, err
	}

	auth := Auth{
		Token:        t,
		RefreshToken: rt,
	}

	return auth, nil
}
