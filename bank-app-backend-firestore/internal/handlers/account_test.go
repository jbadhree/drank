package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock account service
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountService) GetAccountByID(id uint) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountService) GetAccountsByUserID(userID uint) ([]models.Account, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountService) GetAllAccounts() ([]models.Account, error) {
	args := m.Called()
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountService) UpdateAccount(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountService) DeleteAccount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccountService) GenerateAccountNumber() string {
	args := m.Called()
	return args.String(0)
}

func TestGetAllAccounts_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Create test accounts
	testAccounts := []models.Account{
		{ID: 1, UserID: 1, AccountNumber: "1234567890", Balance: 100.0, AccountType: models.Checking},
		{ID: 2, UserID: 2, AccountNumber: "0987654321", Balance: 500.0, AccountType: models.Savings},
	}

	// Set up expectations
	mockAccountService.On("GetAllAccounts").Return(testAccounts, nil)

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	accountHandler.GetAllAccounts(c)

	// Parse the response
	var response []models.AccountDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, response, 2)
	assert.Equal(t, testAccounts[0].ID, response[0].ID)
	assert.Equal(t, testAccounts[1].AccountNumber, response[1].AccountNumber)
	mockAccountService.AssertExpectations(t)
}

func TestGetAllAccounts_Error(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Set up expectations
	mockAccountService.On("GetAllAccounts").Return([]models.Account{}, errors.New("database error"))

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the handler
	accountHandler.GetAllAccounts(c)

	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Failed to get accounts: database error", response.Message)
	mockAccountService.AssertExpectations(t)
}

func TestGetAccountByID_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Create test account
	testAccount := &models.Account{
		ID:            1,
		UserID:        1,
		AccountNumber: "1234567890",
		Balance:       100.0,
		AccountType:   models.Checking,
	}

	// Set up expectations
	mockAccountService.On("GetAccountByID", uint(1)).Return(testAccount, nil)

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/1", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	// Call the handler
	accountHandler.GetAccountByID(c)

	// Parse the response
	var response models.AccountDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testAccount.ID, response.ID)
	assert.Equal(t, testAccount.AccountNumber, response.AccountNumber)
	assert.Equal(t, testAccount.Balance, response.Balance)
	mockAccountService.AssertExpectations(t)
}

func TestGetAccountByID_InvalidID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/invalid", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "invalid"},
	}

	// Call the handler
	accountHandler.GetAccountByID(c)

	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid ID format", response.Message)
	mockAccountService.AssertNotCalled(t, "GetAccountByID")
}

func TestGetAccountByID_NotFound(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Set up expectations for account not found
	mockAccountService.On("GetAccountByID", uint(999)).Return(nil, errors.New("account not found"))

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/999", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "999"},
	}

	// Call the handler
	accountHandler.GetAccountByID(c)

	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Account not found: account not found", response.Message)
	mockAccountService.AssertExpectations(t)
}

func TestGetAccountsByUserID_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Create test accounts
	testAccounts := []models.Account{
		{ID: 1, UserID: 1, AccountNumber: "1234567890", Balance: 100.0, AccountType: models.Checking},
		{ID: 2, UserID: 1, AccountNumber: "0987654321", Balance: 500.0, AccountType: models.Savings},
	}

	// Set up expectations
	mockAccountService.On("GetAccountsByUserID", uint(1)).Return(testAccounts, nil)

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/user/1", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "userId", Value: "1"},
	}

	// Call the handler
	accountHandler.GetAccountsByUserID(c)

	// Parse the response
	var response []models.AccountDTO
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, response, 2)
	assert.Equal(t, testAccounts[0].ID, response[0].ID)
	assert.Equal(t, testAccounts[1].AccountNumber, response[1].AccountNumber)
	mockAccountService.AssertExpectations(t)
}

func TestGetAccountsByUserID_InvalidID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/user/invalid", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "userId", Value: "invalid"},
	}

	// Call the handler
	accountHandler.GetAccountsByUserID(c)

	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid user ID format", response.Message)
	mockAccountService.AssertNotCalled(t, "GetAccountsByUserID")
}

func TestGetAccountsByUserID_Error(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock service
	mockAccountService := new(MockAccountService)

	// Set up expectations for error
	mockAccountService.On("GetAccountsByUserID", uint(999)).Return([]models.Account{}, errors.New("database error"))

	// Create account handler with mock service
	accountHandler := NewAccountHandler(mockAccountService)

	// Create a request to pass to our handler
	req, _ := http.NewRequest("GET", "/api/v1/accounts/user/999", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "userId", Value: "999"},
	}

	// Call the handler
	accountHandler.GetAccountsByUserID(c)

	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert expectations
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Failed to get accounts: database error", response.Message)
	mockAccountService.AssertExpectations(t)
}
