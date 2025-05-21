package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.UserDTO
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get users: " + err.Error()})
		return
	}

	// Convert to DTOs
	userDTOs := make([]models.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToDTO()
	}

	c.JSON(http.StatusOK, userDTOs)
}

// @Summary Get user by ID
// @Description Get a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	user, err := h.userService.GetUserByID(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "User not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.ToDTO())
}

// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserDTO
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.ToDTO())
}
