package unit

import (
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepository interface
var _ interfaces.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAccountRepository implements the AccountRepository interface for testing
type MockAccountRepository struct {
	mock.Mock
}

// Ensure MockAccountRepository implements AccountRepository interface
var _ interfaces.AccountRepository = (*MockAccountRepository)(nil)

func (m *MockAccountRepository) Create(account models.Account) (models.Account, error) {
	args := m.Called(account)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByID(id string) (models.Account, error) {
	args := m.Called(id)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByUserID(userID string) ([]models.Account, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByAccountNumber(accountNumber string) (models.Account, error) {
	args := m.Called(accountNumber)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindAll() ([]models.Account, error) {
	args := m.Called()
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(account models.Account) (models.Account, error) {
	args := m.Called(account)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateBalance(id string, amount float64) (models.Account, error) {
	args := m.Called(id, amount)
	return args.Get(0).(models.Account), args.Error(1)
}

// MockTransactionRepository implements the TransactionRepository interface for testing
type MockTransactionRepository struct {
	mock.Mock
}

// Ensure MockTransactionRepository implements TransactionRepository interface
var _ interfaces.TransactionRepository = (*MockTransactionRepository)(nil)

func (m *MockTransactionRepository) Create(transaction models.Transaction) (models.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByID(id string) (models.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByAccountID(accountID string) ([]models.Transaction, error) {
	args := m.Called(accountID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindAll() ([]models.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindBySourceAccountID(sourceAccountID string) ([]models.Transaction, error) {
	args := m.Called(sourceAccountID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByTargetAccountID(targetAccountID string) ([]models.Transaction, error) {
	args := m.Called(targetAccountID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CreateTransfer(sourceAccountID, targetAccountID string, amount float64, description string) error {
	args := m.Called(sourceAccountID, targetAccountID, amount, description)
	return args.Error(0)
}
