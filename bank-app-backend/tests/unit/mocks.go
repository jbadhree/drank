package unit

import (
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock for UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock for AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByID(id uint) (*models.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByUserID(userID uint) ([]models.Account, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByAccountNumber(accountNumber string) (*models.Account, error) {
	args := m.Called(accountNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) FindAll() ([]models.Account, error) {
	args := m.Called()
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(account *models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccountRepository) UpdateBalance(id uint, amount float64) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByIDWithLock(id uint) (*models.Account, *gorm.DB, error) {
	args := m.Called(id)
	var account *models.Account
	
	if args.Get(0) != nil {
		account = args.Get(0).(*models.Account)
	}
	
	// Return nil for the DB, as we'll use our MockDB for testing
	return account, nil, args.Error(2)
}

// Mock for TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateWithTx(transaction *models.Transaction, tx *gorm.DB) error {
	args := m.Called(transaction, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(id uint) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error) {
	args := m.Called(accountID, limit, offset)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindAll(limit, offset int) ([]models.Transaction, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CountByAccountID(accountID uint) (int64, error) {
	args := m.Called(accountID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTransactionRepository) CountAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockDB for transaction testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() *MockDB {
	m.Called()
	return m
}

func (m *MockDB) Rollback() *MockDB {
	m.Called()
	return m
}

func (m *MockDB) Commit() *MockDB {
	m.Called()
	return m
}

func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Save(value interface{}) *MockDB {
	m.Called(value)
	return m
}

func (m *MockDB) Set(name string, value interface{}) *MockDB {
	m.Called(name, value)
	return m
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *MockDB {
	m.Called(dest, conds)
	return m
}
