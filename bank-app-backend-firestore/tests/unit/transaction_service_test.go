package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService(t *testing.T) {
	// Set up common test data
	now := time.Now()

	t.Run("Create should create a new deposit transaction", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		account := models.Account{
			ID:            "acc123",
			UserID:        "user123",
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		transaction := models.Transaction{
			AccountID:       "acc123",
			Amount:          500.00,
			Type:            models.Deposit,
			Description:     "Test deposit",
			TransactionDate: now,
		}

		updatedAccount := account
		updatedAccount.Balance = 1500.00
		updatedAccount.UpdatedAt = now

		createdTransaction := transaction
		createdTransaction.ID = "t123"
		createdTransaction.Balance = 1500.00
		createdTransaction.CreatedAt = now
		createdTransaction.UpdatedAt = now

		mockAccountRepo.On("FindByID", "acc123").Return(account, nil)
		mockTransactionRepo.On("Create", mock.AnythingOfType("models.Transaction")).Return(createdTransaction, nil)
		mockAccountRepo.On("Update", mock.AnythingOfType("models.Account")).Return(updatedAccount, nil)

		// Act
		result, err := service.Create(transaction)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "t123", result.ID)
		assert.Equal(t, 500.00, result.Amount)
		assert.Equal(t, 1500.00, result.Balance)
		mockTransactionRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Create should create a new withdrawal transaction", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		account := models.Account{
			ID:            "acc123",
			UserID:        "user123",
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		transaction := models.Transaction{
			AccountID:       "acc123",
			Amount:          -200.00,
			Type:            models.Withdrawal,
			Description:     "Test withdrawal",
			TransactionDate: now,
		}

		updatedAccount := account
		updatedAccount.Balance = 800.00
		updatedAccount.UpdatedAt = now

		createdTransaction := transaction
		createdTransaction.ID = "t123"
		createdTransaction.Balance = 800.00
		createdTransaction.CreatedAt = now
		createdTransaction.UpdatedAt = now

		mockAccountRepo.On("FindByID", "acc123").Return(account, nil)
		mockTransactionRepo.On("Create", mock.AnythingOfType("models.Transaction")).Return(createdTransaction, nil)
		mockAccountRepo.On("Update", mock.AnythingOfType("models.Account")).Return(updatedAccount, nil)

		// Act
		result, err := service.Create(transaction)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "t123", result.ID)
		assert.Equal(t, -200.00, result.Amount)
		assert.Equal(t, 800.00, result.Balance)
		mockTransactionRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Transfer should transfer funds between accounts", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		req := models.TransferRequest{
			FromAccountID: "acc123",
			ToAccountID:   "acc456",
			Amount:        200.00,
			Description:   "Test transfer",
		}

		fromAccount := models.Account{
			ID:            "acc123",
			UserID:        "user123",
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		toAccount := models.Account{
			ID:            "acc456",
			UserID:        "user456",
			AccountNumber: "1000000002",
			AccountType:   models.Savings,
			Balance:       500.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		mockAccountRepo.On("FindByID", "acc123").Return(fromAccount, nil)
		mockAccountRepo.On("FindByID", "acc456").Return(toAccount, nil)
		mockTransactionRepo.On("CreateTransfer", "acc123", "acc456", 200.00, "Test transfer").Return(nil)

		// Act
		err := service.Transfer(req)

		// Assert
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
	})

	t.Run("Transfer should fail if source account not found", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		req := models.TransferRequest{
			FromAccountID: "nonexistent",
			ToAccountID:   "acc456",
			Amount:        200.00,
			Description:   "Test transfer",
		}

		mockAccountRepo.On("FindByID", "nonexistent").Return(models.Account{}, errors.New("account not found"))

		// Act
		err := service.Transfer(req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "source account not found")
		mockTransactionRepo.AssertNotCalled(t, "CreateTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Transfer should fail if target account not found", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		req := models.TransferRequest{
			FromAccountID: "acc123",
			ToAccountID:   "nonexistent",
			Amount:        200.00,
			Description:   "Test transfer",
		}

		fromAccount := models.Account{
			ID:            "acc123",
			UserID:        "user123",
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		mockAccountRepo.On("FindByID", "acc123").Return(fromAccount, nil)
		mockAccountRepo.On("FindByID", "nonexistent").Return(models.Account{}, errors.New("account not found"))

		// Act
		err := service.Transfer(req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target account not found")
		mockTransactionRepo.AssertNotCalled(t, "CreateTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Transfer should fail if amount is negative or zero", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		req := models.TransferRequest{
			FromAccountID: "acc123",
			ToAccountID:   "acc456",
			Amount:        0.00,
			Description:   "Test transfer",
		}

		// Act
		err := service.Transfer(req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transfer amount must be positive")
		mockAccountRepo.AssertNotCalled(t, "FindByID", mock.Anything)
		mockTransactionRepo.AssertNotCalled(t, "CreateTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("Transfer should fail if source and target accounts are the same", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		req := models.TransferRequest{
			FromAccountID: "acc123",
			ToAccountID:   "acc123",
			Amount:        200.00,
			Description:   "Test transfer",
		}

		// Act
		err := service.Transfer(req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transfer to the same account")
		mockAccountRepo.AssertNotCalled(t, "FindByID", mock.Anything)
		mockTransactionRepo.AssertNotCalled(t, "CreateTransfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("GetByID should return a transaction when found", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		transaction := models.Transaction{
			ID:              "t123",
			AccountID:       "acc123",
			Amount:          500.00,
			Balance:         1500.00,
			Type:            models.Deposit,
			Description:     "Test deposit",
			TransactionDate: now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		mockTransactionRepo.On("FindByID", "t123").Return(transaction, nil)

		// Act
		result, err := service.GetByID("t123")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "t123", result.ID)
		assert.Equal(t, "acc123", result.AccountID)
		assert.Equal(t, 500.00, result.Amount)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("GetByID should return error when transaction not found", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		mockTransactionRepo.On("FindByID", "nonexistent").Return(models.Transaction{}, errors.New("transaction not found"))

		// Act
		result, err := service.GetByID("nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.TransactionDTO{}, result)
		assert.Contains(t, err.Error(), "transaction not found")
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("GetByAccountID should return transactions for an account", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		transactions := []models.Transaction{
			{
				ID:              "t123",
				AccountID:       "acc123",
				Amount:          500.00,
				Balance:         1500.00,
				Type:            models.Deposit,
				Description:     "Test deposit",
				TransactionDate: now,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				ID:              "t124",
				AccountID:       "acc123",
				Amount:          -200.00,
				Balance:         1300.00,
				Type:            models.Withdrawal,
				Description:     "Test withdrawal",
				TransactionDate: now.Add(-time.Hour),
				CreatedAt:       now.Add(-time.Hour),
				UpdatedAt:       now.Add(-time.Hour),
			},
		}

		mockTransactionRepo.On("FindByAccountID", "acc123").Return(transactions, nil)

		// Act
		result, err := service.GetByAccountID("acc123")

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "t123", result[0].ID)
		assert.Equal(t, "t124", result[1].ID)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("GetAll should return all transactions", func(t *testing.T) {
		// Arrange
		mockTransactionRepo := new(MockTransactionRepository)
		mockAccountRepo := new(MockAccountRepository)
		service := services.NewTransactionService(mockTransactionRepo, mockAccountRepo)

		transactions := []models.Transaction{
			{
				ID:              "t123",
				AccountID:       "acc123",
				Amount:          500.00,
				Balance:         1500.00,
				Type:            models.Deposit,
				Description:     "Test deposit",
				TransactionDate: now,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			{
				ID:              "t124",
				AccountID:       "acc456",
				Amount:          300.00,
				Balance:         800.00,
				Type:            models.Deposit,
				Description:     "Test deposit",
				TransactionDate: now.Add(-time.Hour),
				CreatedAt:       now.Add(-time.Hour),
				UpdatedAt:       now.Add(-time.Hour),
			},
		}

		mockTransactionRepo.On("FindAll").Return(transactions, nil)

		// Act
		result, err := service.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "t123", result[0].ID)
		assert.Equal(t, "t124", result[1].ID)
		mockTransactionRepo.AssertExpectations(t)
	})
}
