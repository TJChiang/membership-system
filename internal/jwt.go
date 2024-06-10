package internal

import (
	"errors"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/golang-jwt/jwt"
	"os"
)

func ParseAndValidateAccessToken(tokenString string) (*jwt.Token, jwt.Claims, error) {
	appSecret := os.Getenv("APP_SECRET")
	token, err := jwt.ParseWithClaims(
		tokenString,
		&generates.JWTAccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(appSecret), nil
		},
	)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*generates.JWTAccessClaims)
	if !ok || !token.Valid {
		return token, nil, errors.New("invalid token")
	}

	return token, claims, nil
}
