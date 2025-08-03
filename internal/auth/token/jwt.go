package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtService interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(encodedToken string) (*Claims, error)
}

type jwtService struct {
	secretkey string
	issure    string
}

func NewJwtService() JwtService {
	return &jwtService{
		secretkey: "secret-key",
		issure:    "gobid",
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(userId string) (string, error) {
	claim := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    s.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.secretkey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.secretkey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
