package internal

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zero-shubham/authsvc/config"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Scopes:   scopes,
		OnBehalf: obo,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authsvc",
			Subject:   subId,
			Audience:  jwt.ClaimStrings{"surveyx_api_gateway"},
			ExpiresAt: jwt.NewNumericDate(now.Add(JwtExp * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	})
	token.Header["alg"] = "HS256"
	token.Header["kid"] = config.JWK_VALUE

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type TokenClaims struct {
	Scopes   []string `json:"scopes"`
	OnBehalf string   `json:"obo"`
	jwt.RegisteredClaims
}

func ParseToken(tokenString string) (TokenClaims, error) {
	var claims TokenClaims
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSignKey), nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	if !token.Valid {
		return TokenClaims{}, errors.New("error while parsing token")
	}

	log.Info().Interface("claims", claims).Msg("valid token")
	return claims, nil
}
