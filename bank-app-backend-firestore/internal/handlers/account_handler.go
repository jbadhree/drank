package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
)

// AccountHandler - Handler for account operations
type AccountHandler struct {
	accountService *services.AccountService
}

// NewAccountHandler - Create a new account handler
func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// GetAllAccounts - Get all accounts endpoint
// @Summary Get all accounts
// @Description Get a list of all accounts
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AccountDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [get]
func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	accounts, err := h.accountService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// GetAccountByID - Get account by ID endpoint
// @Summary Get account by ID
// @Description Get an account by its ID
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} models.AccountDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	id := c.Param("id")

	account, err := h.accountService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountsByUserID - Get accounts by user ID endpoint
// @Summary Get accounts by user ID
// @Description Get accounts for a specific user
// @Tags accounts
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {array} models.AccountDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/user/{userId} [get]
func (h *AccountHandler) GetAccountsByUserID(c *gin.Context) {
	userID := c.Param("userId")

	// Check if the user is requesting their own accounts
	currentUserID, exists := c.Get("userId")
	if !exists || currentUserID.(string) != userID {
		// For security, only allow admin to see other users' accounts
		// In a real app, this would be handled by roles/permissions
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	accounts, err := h.accountService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// CreateAccount - Create account endpoint
// @Summary Create a new account
// @Description Create a new account for a user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account body models.Account true "Account details"
// @Success 201 {object} models.AccountDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var account models.Account

	// Bind the request
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID from the authenticated user
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	account.UserID = userID.(string)

	// Create the account
	createdAccount, err := h.accountService.Create(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAccount)
}
