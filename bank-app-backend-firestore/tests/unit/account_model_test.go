package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountModel(t *testing.T) {
	t.Run("Account.ToDTO should convert Account to AccountDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		account := models.Account{
			ID:            "acc123",
			UserID:        "user123",
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Act
		dto := account.ToDTO()

		// Assert
		assert.Equal(t, account.ID, dto.ID)
		assert.Equal(t, account.UserID, dto.UserID)
		assert.Equal(t, account.AccountNumber, dto.AccountNumber)
		assert.Equal(t, account.AccountType, dto.AccountType)
		assert.Equal(t, account.Balance, dto.Balance)
		assert.Equal(t, account.CreatedAt, dto.CreatedAt)
		assert.Equal(t, account.UpdatedAt, dto.UpdatedAt)
	})

	t.Run("AccountType constants should have correct values", func(t *testing.T) {
		// Assert
		assert.Equal(t, models.AccountType("CHECKING"), models.Checking)
		assert.Equal(t, models.AccountType("SAVINGS"), models.Savings)
	})
}
