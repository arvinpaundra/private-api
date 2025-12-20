package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ Tokenable = (*JWT)(nil)

type JWT struct {
	secret []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: []byte(secret),
	}
}

func (t *JWT) Encode(identifier string, expiry time.Duration, validAfter time.Duration) (string, error) {
	claims := Claims{
		Identifier: identifier,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry).UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(validAfter).UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(t.secret)
}

func (t *JWT) Decode(str string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(str, &Claims{}, func(tkn *jwt.Token) (any, error) {
		_, ok := tkn.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return t.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
