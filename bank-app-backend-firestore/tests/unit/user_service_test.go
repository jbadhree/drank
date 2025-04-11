package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService(t *testing.T) {
	// Set up common test data
	now := time.Now()
	
	t.Run("Create should create a new user", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		inputUser := models.User{
			Email:     "new@example.com",
			Password:  "password123",
			FirstName: "New",
			LastName:  "User",
		}
		
		createdUser := models.User{
			ID:        "user1",
			Email:     "new@example.com",
			Password:  "hashedpassword", // This would be hashed
			FirstName: "New",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		mockRepo.On("Create", mock.AnythingOfType("models.User")).Return(createdUser, nil)
		
		// Act
		userDTO, err := service.Create(inputUser)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user1", userDTO.ID)
		assert.Equal(t, "new@example.com", userDTO.Email)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("GetByID should return user when found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := models.User{
			ID:        "user1",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		mockRepo.On("FindByID", "user1").Return(user, nil)
		
		// Act
		result, err := service.GetByID("user1")
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user1", result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("GetByID should return error when user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		mockRepo.On("FindByID", "nonexistent").Return(models.User{}, errors.New("user not found"))
		
		// Act
		result, err := service.GetByID("nonexistent")
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.UserDTO{}, result)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Authenticate should return user when credentials are valid", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		user := models.User{
			ID:        "user1",
			Email:     "test@example.com",
			Password:  string(hashedPassword),
			FirstName: "Test",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		
		// Act
		result, err := service.Authenticate("test@example.com", password)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user1", result.ID)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Authenticate should return error when email not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		mockRepo.On("FindByEmail", "nonexistent@example.com").Return(models.User{}, errors.New("user not found"))
		
		// Act
		result, err := service.Authenticate("nonexistent@example.com", "password123")
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Authenticate should return error when password is incorrect", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		user := models.User{
			ID:        "user1",
			Email:     "test@example.com",
			Password:  string(hashedPassword),
			FirstName: "Test",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		
		// Act
		result, err := service.Authenticate("test@example.com", "wrongpassword")
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("GetAll should return all users", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		users := []models.User{
			{
				ID:        "user1",
				Email:     "user1@example.com",
				FirstName: "User",
				LastName:  "One",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "user2",
				Email:     "user2@example.com",
				FirstName: "User",
				LastName:  "Two",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}
		
		mockRepo.On("FindAll").Return(users, nil)
		
		// Act
		result, err := service.GetAll()
		
		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "user1", result[0].ID)
		assert.Equal(t, "user2", result[1].ID)
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Update should update a user", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		service := services.NewUserService(mockRepo)
		
		user := models.User{
			ID:        "user1",
			Email:     "update@example.com",
			FirstName: "Updated",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		updatedUser := models.User{
			ID:        "user1",
			Email:     "update@example.com",
			FirstName: "Updated",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}
		
		mockRepo.On("Update", user).Return(updatedUser, nil)
		
		// Act
		result, err := service.Update(user)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "user1", result.ID)
		assert.Equal(t, "Updated", result.FirstName)
		mockRepo.AssertExpectations(t)
	})
}
