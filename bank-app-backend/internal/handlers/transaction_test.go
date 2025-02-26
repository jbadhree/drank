package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock transaction service
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransactionsByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error) {
	args := m.Called(accountID, limit, offset)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionService) Transfer(request *models.TransferRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func TestGetAllTransactions_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create test transactions
	testTransactions := []models.Transaction{
		{ID: 1, AccountID: 1, Amount: 100.0, Balance: 100.0, Type: models.Deposit, TransactionDate: time.Now()},
		{ID: 2, AccountID: 1, Amount: 50.0, Balance: 50.0, Type: models.Withdrawal, TransactionDate: time.Now()},
	}
	
	// Set up expectations
	mockTransactionService.On("GetAllTransactions", 10, 0).Return(testTransactions, nil)
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/transactions?limit=10&offset=0", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// Add query parameters for pagination
	c.Request.URL.RawQuery = "limit=10&offset=0"
	
	// Call the handler
	transactionHandler.GetAllTransactions(c)
	
	// Parse the response
	var response []models.TransactionDTO
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, response, 2)
	assert.Equal(t, testTransactions[0].ID, response[0].ID)
	assert.Equal(t, testTransactions[1].Type, response[1].Type)
	mockTransactionService.AssertExpectations(t)
}

func TestGetAllTransactions_Error(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Set up expectations for error
	mockTransactionService.On("GetAllTransactions", 10, 0).Return([]models.Transaction{}, errors.New("database error"))
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/transactions?limit=10&offset=0", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// Add query parameters for pagination
	c.Request.URL.RawQuery = "limit=10&offset=0"
	
	// Call the handler
	transactionHandler.GetAllTransactions(c)
	
	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "database error", response.Message)
	mockTransactionService.AssertExpectations(t)
}

func TestGetTransactionByID_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create test transaction
	testTransaction := &models.Transaction{
		ID:              1,
		AccountID:       1,
		Amount:          100.0,
		Balance:         100.0,
		Type:            models.Deposit,
		Description:     "Test deposit",
		TransactionDate: time.Now(),
	}
	
	// Set up expectations
	mockTransactionService.On("GetTransactionByID", uint(1)).Return(testTransaction, nil)
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/transactions/1", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}
	
	// Call the handler
	transactionHandler.GetTransactionByID(c)
	
	// Parse the response
	var response models.TransactionDTO
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testTransaction.ID, response.ID)
	assert.Equal(t, testTransaction.Amount, response.Amount)
	assert.Equal(t, testTransaction.Type, response.Type)
	mockTransactionService.AssertExpectations(t)
}

func TestGetTransactionByID_InvalidID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/transactions/invalid", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "invalid"},
	}
	
	// Call the handler
	transactionHandler.GetTransactionByID(c)
	
	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, response.Message, "invalid transaction ID")
	mockTransactionService.AssertNotCalled(t, "GetTransactionByID")
}

func TestGetTransactionsByAccountID_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create test transactions
	testTransactions := []models.Transaction{
		{ID: 1, AccountID: 1, Amount: 100.0, Balance: 100.0, Type: models.Deposit, TransactionDate: time.Now()},
		{ID: 2, AccountID: 1, Amount: 50.0, Balance: 50.0, Type: models.Withdrawal, TransactionDate: time.Now()},
	}
	
	// Set up expectations
	mockTransactionService.On("GetTransactionsByAccountID", uint(1), 10, 0).Return(testTransactions, nil)
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/transactions/account/1?limit=10&offset=0", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "accountId", Value: "1"},
	}
	// Add query parameters for pagination
	c.Request.URL.RawQuery = "limit=10&offset=0"
	
	// Call the handler
	transactionHandler.GetTransactionsByAccountID(c)
	
	// Parse the response
	var response []models.TransactionDTO
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, response, 2)
	assert.Equal(t, testTransactions[0].ID, response[0].ID)
	assert.Equal(t, testTransactions[1].Type, response[1].Type)
	mockTransactionService.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create transfer request
	transferRequest := models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        50.0,
		Description:   "Test transfer",
	}
	
	// Set up expectations
	mockTransactionService.On("Transfer", mock.MatchedBy(func(req *models.TransferRequest) bool {
		return req.FromAccountID == 1 && req.ToAccountID == 2 && req.Amount == 50.0
	})).Return(nil)
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request body
	jsonValue, _ := json.Marshal(transferRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/transactions/transfer", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	transactionHandler.Transfer(c)
	
	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	mockTransactionService.AssertExpectations(t)
}

func TestTransfer_InvalidRequest(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create invalid transfer request (missing required fields)
	invalidRequest := struct {
		FromAccountID uint `json:"fromAccountId"`
	}{
		FromAccountID: 1,
	}
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request body
	jsonValue, _ := json.Marshal(invalidRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/transactions/transfer", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	transactionHandler.Transfer(c)
	
	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockTransactionService.AssertNotCalled(t, "Transfer")
}

func TestTransfer_ServiceError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockTransactionService := new(MockTransactionService)
	
	// Create transfer request
	transferRequest := models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        50.0,
		Description:   "Test transfer",
	}
	
	// Set up expectations for error
	mockTransactionService.On("Transfer", mock.MatchedBy(func(req *models.TransferRequest) bool {
		return req.FromAccountID == 1 && req.ToAccountID == 2 && req.Amount == 50.0
	})).Return(errors.New("insufficient funds"))
	
	// Create transaction handler with mock service
	transactionHandler := NewTransactionHandler(mockTransactionService)
	
	// Create a request body
	jsonValue, _ := json.Marshal(transferRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/transactions/transfer", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	transactionHandler.Transfer(c)
	
	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "insufficient funds", response.Message)
	mockTransactionService.AssertExpectations(t)
}
