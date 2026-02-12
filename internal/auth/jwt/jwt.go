package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

type Claims struct {
	jwt.RegisteredClaims
	UserId int32 `json:"user_id"`
}

func NewJwtManager(secret string, ttl time.Duration) *JwtManager {
	return &JwtManager{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: "chatter",
	}
}

func (m *JwtManager) Generate(userId int32) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(m.ttl)

	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, expires, nil
}

func (m *JwtManager) GetClaims(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
