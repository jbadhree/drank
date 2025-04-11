package interfaces

import (
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
)

// TransactionRepository defines the interface for transaction repository operations
type TransactionRepository interface {
	Create(transaction models.Transaction) (models.Transaction, error)
	FindByID(id string) (models.Transaction, error)
	FindByAccountID(accountID string) ([]models.Transaction, error)
	FindAll() ([]models.Transaction, error)
	FindBySourceAccountID(sourceAccountID string) ([]models.Transaction, error)
	FindByTargetAccountID(targetAccountID string) ([]models.Transaction, error)
	CreateTransfer(sourceAccountID, targetAccountID string, amount float64, description string) error
}
