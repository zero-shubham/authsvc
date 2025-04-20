package internal

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	JwtSignEnv = "JWT_SIGN_KEY"
	JwtExp     = 900
)

var jwtSignKey = os.Getenv(JwtSignEnv)

func CreateToken(subId string, scopes []string, obo string) (string, error) {
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
		"jti":    uuid.New().String(),
		"obo":    obo,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type TokenClaims struct {
	Audience string   `json:"aud"`
	Issuer   string   `json:"iss"`
	Subject  string   `json:"sub"`
	Foo      string   `json:"foo"`
	IssuedAt int64    `json:"iat"`
	Expires  int64    `json:"exp"`
	Scopes   []string `json:"scopes"`
	JwtID    string   `json:"jti"`
	OBO      string   `json:"obo"`
}

func ParseToken(tokenString string) (TokenClaims, error) {
	claims := jwt.MapClaims{}
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSignKey, nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	if !token.Valid {
		return TokenClaims{}, errors.New("error while parsing token")
	}

	return TokenClaims{
		Audience: claims["aud"].(string),
		Issuer:   claims["iss"].(string),
		Subject:  claims["sub"].(string),
		Foo:      claims["foo"].(string),
		IssuedAt: int64(claims["iat"].(float64)),
		Expires:  int64(claims["exp"].(float64)),
		Scopes:   claims["scopes"].([]string),
		JwtID:    claims["jti"].(string),
		OBO:      claims["obo"].(string),
	}, nil
}
