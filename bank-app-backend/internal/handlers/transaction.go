package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService}
}

// @Summary Get all transactions
// @Description Get a paginated list of all transactions
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.TransactionDTO
// @Failure 500 {object} ErrorResponse
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.transactionService.GetAllTransactions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get transactions: " + err.Error()})
		return
	}

	// Convert to DTOs
	transactionDTOs := make([]models.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		transactionDTOs[i] = transaction.ToDTO()
	}

	c.JSON(http.StatusOK, transactionDTOs)
}

// @Summary Get transaction by ID
// @Description Get a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.TransactionDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid ID format"})
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Transaction not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction.ToDTO())
}

// @Summary Get transactions by account ID
// @Description Get a paginated list of transactions for an account
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param accountId path int true "Account ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.TransactionDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/account/{accountId} [get]
func (h *TransactionHandler) GetTransactionsByAccountID(c *gin.Context) {
	accountIDStr := c.Param("accountId")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid account ID format"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.transactionService.GetTransactionsByAccountID(uint(accountID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to get transactions: " + err.Error()})
		return
	}

	// Convert to DTOs
	transactionDTOs := make([]models.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		transactionDTOs[i] = transaction.ToDTO()
	}

	c.JSON(http.StatusOK, transactionDTOs)
}

// @Summary Transfer money
// @Description Transfer money between accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transferRequest body models.TransferRequest true "Transfer Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /transactions/transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	var transferRequest models.TransferRequest
	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request: " + err.Error()})
		return
	}

	if err := h.transactionService.Transfer(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Transfer failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}
