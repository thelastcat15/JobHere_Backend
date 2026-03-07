package utils

import (
	"errors"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

var jwks *keyfunc.JWKS

// InitJWKS เรียกตอน server start
func InitJWKS(jwksURL string) error {
	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			// log error ได้ตรงนี้
		},
	}

	var err error
	jwks, err = keyfunc.Get(jwksURL, options)
	return err
}

type SupabaseClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func ParseSupabaseJWT(tokenString string) (*SupabaseClaims, error) {
	if jwks == nil {
		return nil, errors.New("JWKS not initialized")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&SupabaseClaims{},
		jwks.Keyfunc,
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SupabaseClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
