package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
)

type JwtBuilder struct {
	secretKey string
}

func NewJwtBuilder(secretKey string) *JwtBuilder {
	return &JwtBuilder{secretKey}
}

func (builder JwtBuilder) GenerateToken(id uint, email string, duration time.Duration) (string, *UserClaims, error) {
	var userClaims *UserClaims
	userClaims, err := NewUserClaims(id, email, duration)
	if err != nil {
		return "", nil, err
	}

	var tokenJwt *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	token, err := tokenJwt.SignedString([]byte(builder.secretKey))
	if err != nil {
		return "", nil, err
	}

	return token, userClaims, nil
}

func (builder JwtBuilder) VerifyToken(token string) (*UserClaims, error) {
	userClaims := &UserClaims{}
	var tokenJwt *jwt.Token
	tokenJwt, err := jwt.ParseWithClaims(token, userClaims, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errorhandler.ErrInvalidTokenSigningMethod
		}
		return []byte(builder.secretKey), nil
	})

	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errorhandler.ErrInvalidExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errorhandler.ErrTokenNotValidYet
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errorhandler.ErrMalformedToken
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errorhandler.ErrInvalidTokenSignature
		}

		return nil, errorhandler.ErrParsingToken
	}

	userClaims, ok := tokenJwt.Claims.(*UserClaims)
	if !ok {
		return nil, errorhandler.ErrInvalidTokenClaim
	}

	return userClaims, nil
}
