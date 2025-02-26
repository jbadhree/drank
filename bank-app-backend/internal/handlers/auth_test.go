package handlers

import (
	"bytes"
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

// Mock user service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) Authenticate(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockUserService := new(MockUserService)
	
	// Create test user
	testUser := &models.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	
	// Set up expectations
	mockUserService.On("Authenticate", "test@example.com", "password123").Return(testUser, nil)
	
	// Create auth handler with mock service
	jwtSecret := "test-secret-key"
	authHandler := NewAuthHandler(mockUserService, jwtSecret)
	
	// Create a request body
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonValue, _ := json.Marshal(loginRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	authHandler.Login(c)
	
	// Parse the response
	var response models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, testUser.ID, response.User.ID)
	assert.Equal(t, testUser.Email, response.User.Email)
	mockUserService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockUserService := new(MockUserService)
	
	// Set up expectations for invalid credentials
	mockUserService.On("Authenticate", "test@example.com", "wrongpassword").
		Return(nil, errors.New("invalid credentials"))
	
	// Create auth handler with mock service
	jwtSecret := "test-secret-key"
	authHandler := NewAuthHandler(mockUserService, jwtSecret)
	
	// Create a request body
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonValue, _ := json.Marshal(loginRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	authHandler.Login(c)
	
	// Parse the response
	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Assert expectations
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "invalid credentials", response.Message)
	mockUserService.AssertExpectations(t)
}

func TestLogin_InvalidRequest(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create mock service
	mockUserService := new(MockUserService)
	
	// Create auth handler with mock service
	jwtSecret := "test-secret-key"
	authHandler := NewAuthHandler(mockUserService, jwtSecret)
	
	// Create an invalid request body (missing required fields)
	loginRequest := struct {
		Email string `json:"email"`
	}{
		Email: "test@example.com",
	}
	jsonValue, _ := json.Marshal(loginRequest)
	
	// Create a request to pass to our handler
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create a gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	// Call the handler
	authHandler.Login(c)
	
	// Assert expectations
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUserService.AssertNotCalled(t, "Authenticate")
}
