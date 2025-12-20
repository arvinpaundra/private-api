package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokenable interface {
	Encode(identifier string, expiry time.Duration, validAfter time.Duration) (string, error)
	Decode(str string) (*Claims, error)
}

type Claims struct {
	Identifier string `json:"identifier"`
	jwt.RegisteredClaims
}
