package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/middleware"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
)

// AuthHandler - Handler for authentication operations
type AuthHandler struct {
	userService *services.UserService
	jwtSecret   string
}

// NewAuthHandler - Create a new auth handler
func NewAuthHandler(userService *services.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// Login - Login endpoint
// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate the user
	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token
	authMiddleware := middleware.NewAuthMiddleware(h.jwtSecret)
	token, err := authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token and user
	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user.ToDTO(),
	})
}

// Register - Register endpoint
// @Summary Register user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body models.User true "User details"
// @Success 201 {object} models.UserDTO
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User

	// Bind the request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the user
	createdUser, err := h.userService.Create(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the created user
	c.JSON(http.StatusCreated, createdUser)
}
