package interfaces

import (
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
)

// AccountRepository defines the interface for account repository operations
type AccountRepository interface {
	Create(account models.Account) (models.Account, error)
	FindByID(id string) (models.Account, error)
	FindByUserID(userID string) ([]models.Account, error)
	FindByAccountNumber(accountNumber string) (models.Account, error)
	FindAll() ([]models.Account, error)
	Update(account models.Account) (models.Account, error)
	Delete(id string) error
	UpdateBalance(id string, amount float64) (models.Account, error)
}
