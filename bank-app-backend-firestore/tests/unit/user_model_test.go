package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	t.Run("BeforeSave should hash password", func(t *testing.T) {
		// Arrange
		user := &models.User{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		// Act
		err := user.BeforeSave(&gorm.DB{})

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, "password123", user.Password)
		// Verify it's a valid bcrypt hash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err)
	})

	t.Run("BeforeSave should not hash empty password", func(t *testing.T) {
		// Arrange
		user := &models.User{
			Email:     "test@example.com",
			Password:  "",
			FirstName: "Test",
			LastName:  "User",
		}

		// Act
		err := user.BeforeSave(&gorm.DB{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", user.Password)
	})

	t.Run("ComparePassword should verify correct password", func(t *testing.T) {
		// Arrange
		password := "password123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &models.User{
			Password: string(hash),
		}

		// Act
		err := user.ComparePassword(password)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("ComparePassword should reject incorrect password", func(t *testing.T) {
		// Arrange
		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &models.User{
			Password: string(hash),
		}

		// Act
		err := user.ComparePassword("wrongpassword")

		// Assert
		assert.Error(t, err)
	})

	t.Run("ToDTO should convert User to UserDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		user := &models.User{
			ID:        1,
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
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "test@example.com", dto.Email)
		assert.Equal(t, "Test", dto.FirstName)
		assert.Equal(t, "User", dto.LastName)
		assert.Equal(t, now, dto.CreatedAt)
		assert.Equal(t, now, dto.UpdatedAt)
	})
}
