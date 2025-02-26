package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountModel(t *testing.T) {
	t.Run("ToDTO should convert Account to AccountDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		account := &models.Account{
			ID:            1,
			UserID:        2,
			AccountNumber: "ACC12345",
			AccountType:   models.Checking,
			Balance:       1000.50,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Act
		dto := account.ToDTO()

		// Assert
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, uint(2), dto.UserID)
		assert.Equal(t, "ACC12345", dto.AccountNumber)
		assert.Equal(t, models.Checking, dto.AccountType)
		assert.Equal(t, 1000.50, dto.Balance)
		assert.Equal(t, now, dto.CreatedAt)
		assert.Equal(t, now, dto.UpdatedAt)
	})

	t.Run("Account types should be defined correctly", func(t *testing.T) {
		// Assert
		assert.Equal(t, models.AccountType("CHECKING"), models.Checking)
		assert.Equal(t, models.AccountType("SAVINGS"), models.Savings)
	})
}
