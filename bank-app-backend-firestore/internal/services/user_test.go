package services

import (
	"errors"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Create a mock for the user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthenticateUser_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Create a test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	
	// Set up expectations
	mockRepo.On("FindByEmail", "test@example.com").Return(testUser, nil)
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method being tested
	user, err := service.AuthenticateUser("test@example.com", "password123")
	
	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.ID, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Create a test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	
	// Set up expectations
	mockRepo.On("FindByEmail", "test@example.com").Return(testUser, nil)
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method with wrong password
	user, err := service.AuthenticateUser("test@example.com", "wrongpassword")
	
	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Set up expectations for user not found
	mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method
	user, err := service.AuthenticateUser("nonexistent@example.com", "password123")
	
	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Create a test user
	testUser := &models.User{
		ID:    1,
		Email: "test@example.com",
		FirstName: "Test",
		LastName: "User",
	}
	
	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(testUser, nil)
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method being tested
	user, err := service.GetUserByID(1)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Set up expectations for user not found
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method
	user, err := service.GetUserByID(999)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockUserRepository)
	
	// Create test users
	testUsers := []models.User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}
	
	// Set up expectations
	mockRepo.On("FindAll").Return(testUsers, nil)
	
	// Create service with mock repo
	service := NewUserService(mockRepo)
	
	// Call the method being tested
	users, err := service.GetAllUsers()
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, testUsers[0].ID, users[0].ID)
	assert.Equal(t, testUsers[1].Email, users[1].Email)
	mockRepo.AssertExpectations(t)
}
