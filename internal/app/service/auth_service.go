package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"propmanager/internal/config"
)

type AuthService struct {
	Cfg *config.AuthConfig
}

func NewAuthService(cfg *config.AuthConfig) *AuthService {
	return &AuthService{Cfg: cfg}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Cfg.SecretKey))
}
