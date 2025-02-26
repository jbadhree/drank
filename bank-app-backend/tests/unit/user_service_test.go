package unit

import (
	"errors"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService(t *testing.T) {
	t.Run("CreateUser should create a new user when email doesn't exist", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := &models.User{
			Email:     "new@example.com",
			Password:  "password123",
			FirstName: "New",
			LastName:  "User",
		}
		
		mockRepo.On("FindByEmail", "new@example.com").Return(nil, errors.New("user not found"))
		mockRepo.On("Create", user).Return(nil)
		
		// Act
		err := service.CreateUser(user)
		
		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("CreateUser should return error when email already exists", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := &models.User{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "Existing",
			LastName:  "User",
		}
		
		existingUser := &models.User{
			ID:        1,
			Email:     "existing@example.com",
			FirstName: "Existing",
			LastName:  "User",
		}
		
		mockRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil)
		
		// Act
		err := service.CreateUser(user)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "user with this email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("GetUserByID should return user when found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}
		
		mockRepo.On("FindByID", uint(1)).Return(user, nil)
		
		// Act
		result, err := service.GetUserByID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("GetUserByID should return error when user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
		
		// Act
		result, err := service.GetUserByID(999)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("AuthenticateUser should return user when credentials are valid", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  string(hashedPassword),
			FirstName: "Test",
			LastName:  "User",
		}
		
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		
		// Act
		result, err := service.AuthenticateUser("test@example.com", password)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("AuthenticateUser should return error when email not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
		
		// Act
		result, err := service.AuthenticateUser("nonexistent@example.com", "password123")
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("AuthenticateUser should return error when password is incorrect", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		user := &models.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  string(hashedPassword),
			FirstName: "Test",
			LastName:  "User",
		}
		
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		
		// Act
		result, err := service.AuthenticateUser("test@example.com", "wrongpassword")
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("UpdateUser should update when user exists and email not taken", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := &models.User{
			ID:        1,
			Email:     "update@example.com",
			FirstName: "Updated",
			LastName:  "User",
		}
		
		mockRepo.On("FindByID", uint(1)).Return(user, nil)
		mockRepo.On("FindByEmail", "update@example.com").Return(nil, errors.New("user not found"))
		mockRepo.On("Update", user).Return(nil)
		
		// Act
		err := service.UpdateUser(user)
		
		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("UpdateUser should fail when email already taken", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := &models.User{
			ID:        1,
			Email:     "taken@example.com",
			FirstName: "Updated",
			LastName:  "User",
		}
		
		existingUser := &models.User{
			ID:        2, // Different ID
			Email:     "taken@example.com",
			FirstName: "Existing",
			LastName:  "User",
		}
		
		mockRepo.On("FindByID", uint(1)).Return(user, nil)
		mockRepo.On("FindByEmail", "taken@example.com").Return(existingUser, nil)
		
		// Act
		err := service.UpdateUser(user)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "email is already taken by another user", err.Error())
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}
