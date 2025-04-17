package internal

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JwtSignEnv = "JWT_SIGN_KEY"
	JwtExp     = 900
)

var jwtSignKey = os.Getenv(JwtSignEnv)

func CreateToken(subId string, scopes []string) (string, error) {
	now := time.Now()
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":    "surveyx_api_gateway",
		"iss":    "authsvc",
		"sub":    subId,
		"foo":    "bar",
		"iat":    now.Unix(),
		"exp":    now.Add(JwtExp * time.Second).Unix(),
		"scopes": scopes,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
