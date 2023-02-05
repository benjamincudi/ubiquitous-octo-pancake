package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	authToken = "atillmtoken"
	jwtSecret = "super_secret_signing_key_abc_123456"
)

type UserInfo struct {
	PersonalAccountNumber int `json:"pan" form:"pan"`
}

func getSignedJwt(_ context.Context, info UserInfo) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"pan": info.PersonalAccountNumber,
			"exp": time.Now().Add(time.Hour).Add(time.Minute).Unix(),
		},
	).SignedString(jwtSecret)
}

func userFromContext(c *gin.Context) (UserInfo, error) {
	tCookie, err := c.Request.Cookie(authToken)
	if err != nil && errors.Is(err, http.ErrNoCookie) {
		log.Println("no auth cookie found")
		return UserInfo{}, err
	}
	var userInfo struct {
		UserInfo
		jwt.RegisteredClaims
	}
	if _, pErr := jwt.ParseWithClaims(tCookie.Value, &userInfo, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}); pErr != nil {
		log.Printf("invalid jwt: %v\n", pErr)
		return UserInfo{}, err
	}
	return userInfo.UserInfo, nil
}
