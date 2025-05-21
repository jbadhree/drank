package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
)

type AccountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{accountService}
}

// @Summary Get all accounts
// @Description Get a list of all accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AccountDTO
// @Failure 500 {object} ErrorResponse
// @Router /accounts [get]
func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get accounts: " + err.Error()})
		return
	}

	// Convert to DTOs
	accountDTOs := make([]models.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}

	c.JSON(http.StatusOK, accountDTOs)
}

// @Summary Get account by ID
// @Description Get an account by its ID
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Account ID"
// @Success 200 {object} models.AccountDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	idStr := c.Param("id")

	account, err := h.accountService.GetAccountByID(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Account not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, account.ToDTO())
}

// @Summary Get accounts by user ID
// @Description Get all accounts for a user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {array} models.AccountDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/user/{userId} [get]
func (h *AccountHandler) GetAccountsByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")

	accounts, err := h.accountService.GetAccountsByUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get accounts: " + err.Error()})
		return
	}

	// Convert to DTOs
	accountDTOs := make([]models.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}

	c.JSON(http.StatusOK, accountDTOs)
}
