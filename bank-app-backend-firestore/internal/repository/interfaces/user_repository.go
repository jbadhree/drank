package interfaces

import (
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByID(id string) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindAll() ([]models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id string) error
}
