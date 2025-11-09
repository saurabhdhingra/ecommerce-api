package service

import (
	"errors"
	"time"
	"net/http"

	"ecommerce-api/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthKey string
const UserContextKey AuthKey = "user"

type JWTService interface {
	GenerateToken(userID uint, isAdmin bool) (string, error)
	ValidateToken(tokenString string) (*domain.Claims, error)
	Middleware(next http.HandlerFunc, requiredAdmin bool) http.HandlerFunc
}

type JWTAuthService struct {
	JWTService
	secret string
}

func NewJWTService(secret string) JWTService {
	return &JWTAuthService{secret: secret}
}

// GenerateToken creates a signed JWT for the given user.
func (s *JWTAuthService) GenerateToken(userID uint, isAdmin bool) (string, error) {
	claims := domain.Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ValidateToken parses and validates the JWT.
func (s *JWTAuthService) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token claims")
}