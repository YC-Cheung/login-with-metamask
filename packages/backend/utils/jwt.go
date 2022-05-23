package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	Uid           uint   `json:"uid"`
	PublicAddress string `json:"publicAddress"`
	jwt.StandardClaims
}

func GenerateToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("metamasklogintest"))

	return tokenString
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("metamasklogintest"), nil
	})

	if err != nil {
		return nil, errors.New("token is invalid")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("couldn't handle this token")
}
