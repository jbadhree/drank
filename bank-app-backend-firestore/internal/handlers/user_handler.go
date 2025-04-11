package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
)

// UserHandler - Handler for user operations
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler - Create a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsers - Get all users endpoint
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.UserDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID - Get user by ID endpoint
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.UserDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUser - Get current user endpoint
// @Summary Get current user
// @Description Get the current authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserDTO
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
