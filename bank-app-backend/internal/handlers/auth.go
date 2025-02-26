package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
)

type AuthHandler struct {
	userService services.UserService
	jwtSecret   string
}

func NewAuthHandler(userService services.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{userService, jwtSecret}
}

// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request: " + err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Authentication failed: " + err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user.ToDTO(),
	})
}

func (h *AuthHandler) generateToken(user *models.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
