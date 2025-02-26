package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactionModel(t *testing.T) {
	t.Run("ToDTO should convert Transaction to TransactionDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		sourceAccID := uint(2)
		targetAccID := uint(3)
		transaction := &models.Transaction{
			ID:              1,
			AccountID:       2,
			SourceAccountID: &sourceAccID,
			TargetAccountID: &targetAccID,
			Amount:          500.75,
			Balance:         1500.25,
			Type:            models.Transfer,
			Description:     "Test transfer",
			TransactionDate: now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		// Act
		dto := transaction.ToDTO()

		// Assert
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, uint(2), dto.AccountID)
		assert.Equal(t, &sourceAccID, dto.SourceAccountID)
		assert.Equal(t, &targetAccID, dto.TargetAccountID)
		assert.Equal(t, 500.75, dto.Amount)
		assert.Equal(t, 1500.25, dto.Balance)
		assert.Equal(t, models.Transfer, dto.Type)
		assert.Equal(t, "Test transfer", dto.Description)
		assert.Equal(t, now, dto.TransactionDate)
		assert.Equal(t, now, dto.CreatedAt)
		assert.Equal(t, now, dto.UpdatedAt)
	})

	t.Run("Transaction types should be defined correctly", func(t *testing.T) {
		// Assert
		assert.Equal(t, models.TransactionType("DEPOSIT"), models.Deposit)
		assert.Equal(t, models.TransactionType("WITHDRAWAL"), models.Withdrawal)
		assert.Equal(t, models.TransactionType("TRANSFER"), models.Transfer)
	})
}
