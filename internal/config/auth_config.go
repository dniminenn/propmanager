package config

import (
	"os"

	"propmanager/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

type AuthConfig struct {
	Username  string
	Password  string
	SecretKey string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		Username:  os.Getenv("AUTH_USERNAME"),
		Password:  os.Getenv("AUTH_PASSWORD"),
		SecretKey: os.Getenv("AUTH_SECRET_KEY"),
	}
}

func (cfg *AuthConfig) AuthMiddleware() gin.HandlerFunc {
	return middleware.AuthMiddleware(cfg.SecretKey)
}
