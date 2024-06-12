package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"propmanager/internal/app/service"
)

// AuthHandler represents the authentication handler.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler returns a new authentication handler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login to the system
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != h.authService.Cfg.Username || password != h.authService.Cfg.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		fmt.Println("Invalid username or password, expected:", h.authService.Cfg.Username, h.authService.Cfg.Password, " got:", username, password)
		return
	}

	token, err := h.authService.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
