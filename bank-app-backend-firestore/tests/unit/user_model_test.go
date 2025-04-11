package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserModel(t *testing.T) {
	t.Run("User.ToDTO should convert User to UserDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		user := models.User{
			ID:        "user123",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			FirstName: "Test",
			LastName:  "User",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		dto := user.ToDTO()

		// Assert
		assert.Equal(t, user.ID, dto.ID)
		assert.Equal(t, user.Email, dto.Email)
		assert.Equal(t, user.FirstName, dto.FirstName)
		assert.Equal(t, user.LastName, dto.LastName)
		assert.Equal(t, user.CreatedAt, dto.CreatedAt)
		assert.Equal(t, user.UpdatedAt, dto.UpdatedAt)
		// Password should not be included in DTO
		assert.Empty(t, "", "DTO should not contain password")
	})

	t.Run("GeneratePasswordHash should hash password", func(t *testing.T) {
		// Arrange
		password := "test123"

		// Act
		hash, err := models.GeneratePasswordHash(password)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, password, hash)
		
		// Verify that the hash can be compared to original password
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		assert.NoError(t, err)
	})

	t.Run("ComparePassword should return nil for correct password", func(t *testing.T) {
		// Arrange
		password := "test123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := models.User{
			Password: string(hash),
		}

		// Act
		err := user.ComparePassword(password)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("ComparePassword should return error for incorrect password", func(t *testing.T) {
		// Arrange
		password := "test123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := models.User{
			Password: string(hash),
		}

		// Act
		err := user.ComparePassword("wrongpassword")

		// Assert
		assert.Error(t, err)
	})
}
