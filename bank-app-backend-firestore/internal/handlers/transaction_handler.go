package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
)

// TransactionHandler - Handler for transaction operations
type TransactionHandler struct {
	transactionService *services.TransactionService
}

// NewTransactionHandler - Create a new transaction handler
func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetAllTransactions - Get all transactions endpoint
// @Summary Get all transactions
// @Description Get a list of all transactions
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.TransactionDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.transactionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID - Get transaction by ID endpoint
// @Summary Get transaction by ID
// @Description Get a transaction by its ID
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.TransactionDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetTransactionsByAccountID - Get transactions by account ID endpoint
// @Summary Get transactions by account ID
// @Description Get transactions for a specific account
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param accountId path string true "Account ID"
// @Success 200 {array} models.TransactionDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/account/{accountId} [get]
func (h *TransactionHandler) GetTransactionsByAccountID(c *gin.Context) {
	accountID := c.Param("accountId")

	transactions, err := h.transactionService.GetByAccountID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Transfer - Transfer funds endpoint
// @Summary Transfer funds
// @Description Transfer funds between accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transferRequest body models.TransferRequest true "Transfer details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /transactions/transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req models.TransferRequest

	// Bind the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform the transfer
	err := h.transactionService.Transfer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

// CreateDeposit - Create deposit endpoint
// @Summary Create a deposit
// @Description Create a deposit to an account
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body models.Transaction true "Deposit details"
// @Success 201 {object} models.TransactionDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /transactions/deposit [post]
func (h *TransactionHandler) CreateDeposit(c *gin.Context) {
	var transaction models.Transaction

	// Bind the request
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transaction type to deposit
	transaction.Type = models.Deposit

	// Create the transaction
	createdTransaction, err := h.transactionService.Create(transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTransaction)
}

// CreateWithdrawal - Create withdrawal endpoint
// @Summary Create a withdrawal
// @Description Create a withdrawal from an account
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body models.Transaction true "Withdrawal details"
// @Success 201 {object} models.TransactionDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /transactions/withdrawal [post]
func (h *TransactionHandler) CreateWithdrawal(c *gin.Context) {
	var transaction models.Transaction

	// Bind the request
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set transaction type to withdrawal
	transaction.Type = models.Withdrawal
	// Make sure amount is negative for withdrawals
	if transaction.Amount > 0 {
		transaction.Amount = -transaction.Amount
	}

	// Create the transaction
	createdTransaction, err := h.transactionService.Create(transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTransaction)
}
