package unit

import (
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactionModel(t *testing.T) {
	t.Run("Transaction.ToDTO should convert Transaction to TransactionDTO", func(t *testing.T) {
		// Arrange
		now := time.Now()
		sourceAccountID := "acc123"
		targetAccountID := "acc456"
		
		transaction := models.Transaction{
			ID:               "t123",
			AccountID:        "acc123",
			SourceAccountID:  &sourceAccountID,
			TargetAccountID:  &targetAccountID,
			Amount:           100.00,
			Balance:          900.00,
			Type:             models.Transfer,
			Description:      "Test transfer",
			TransactionDate:  now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		// Act
		dto := transaction.ToDTO()

		// Assert
		assert.Equal(t, transaction.ID, dto.ID)
		assert.Equal(t, transaction.AccountID, dto.AccountID)
		assert.Equal(t, transaction.SourceAccountID, dto.SourceAccountID)
		assert.Equal(t, transaction.TargetAccountID, dto.TargetAccountID)
		assert.Equal(t, transaction.Amount, dto.Amount)
		assert.Equal(t, transaction.Balance, dto.Balance)
		assert.Equal(t, transaction.Type, dto.Type)
		assert.Equal(t, transaction.Description, dto.Description)
		assert.Equal(t, transaction.TransactionDate, dto.TransactionDate)
		assert.Equal(t, transaction.CreatedAt, dto.CreatedAt)
		assert.Equal(t, transaction.UpdatedAt, dto.UpdatedAt)
	})

	t.Run("TransactionType constants should have correct values", func(t *testing.T) {
		// Assert
		assert.Equal(t, models.TransactionType("DEPOSIT"), models.Deposit)
		assert.Equal(t, models.TransactionType("WITHDRAWAL"), models.Withdrawal)
		assert.Equal(t, models.TransactionType("TRANSFER"), models.Transfer)
	})
}
