package services

import (
	"errors"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Create a mock for the account repository
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
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*models.Account), args.Get(1).(*gorm.DB), args.Error(2)
}

func TestCreateAccount_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account
	testAccount := &models.Account{
		UserID:      1,
		AccountType: models.Checking,
		Balance:     100.0,
	}
	
	// Set up expectations
	mockRepo.On("Create", mock.AnythingOfType("*models.Account")).Return(nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.CreateAccount(testAccount)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.NotEmpty(t, testAccount.AccountNumber) // Account number should be generated
	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_InvalidType(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account with invalid type
	testAccount := &models.Account{
		UserID:      1,
		AccountType: "INVALID_TYPE",
		Balance:     100.0,
	}
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.CreateAccount(testAccount)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "invalid account type", err.Error())
	// Repository method should not be called
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetAccountByID_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account
	testAccount := &models.Account{
		ID:           1,
		UserID:       1,
		AccountType:  models.Checking,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	account, err := service.GetAccountByID(1)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, testAccount.ID, account.ID)
	assert.Equal(t, testAccount.Balance, account.Balance)
	mockRepo.AssertExpectations(t)
}

func TestGetAccountByID_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up expectations for account not found
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("account not found"))
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	account, err := service.GetAccountByID(999)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, "account not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetAccountsByUserID_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test accounts
	testAccounts := []models.Account{
		{ID: 1, UserID: 1, AccountType: models.Checking, Balance: 100.0},
		{ID: 2, UserID: 1, AccountType: models.Savings, Balance: 500.0},
	}
	
	// Set up expectations
	mockRepo.On("FindByUserID", uint(1)).Return(testAccounts, nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	accounts, err := service.GetAccountsByUserID(1)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, testAccounts[0].ID, accounts[0].ID)
	assert.Equal(t, testAccounts[1].Balance, accounts[1].Balance)
	mockRepo.AssertExpectations(t)
}

func TestGetAllAccounts_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test accounts
	testAccounts := []models.Account{
		{ID: 1, UserID: 1, AccountType: models.Checking, Balance: 100.0},
		{ID: 2, UserID: 2, AccountType: models.Savings, Balance: 500.0},
	}
	
	// Set up expectations
	mockRepo.On("FindAll").Return(testAccounts, nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	accounts, err := service.GetAllAccounts()
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, testAccounts[0].ID, accounts[0].ID)
	assert.Equal(t, testAccounts[1].Balance, accounts[1].Balance)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccount_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account
	testAccount := &models.Account{
		ID:           1,
		UserID:       1,
		AccountType:  models.Checking,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	mockRepo.On("Update", testAccount).Return(nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.UpdateAccount(testAccount)
	
	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateAccount_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account
	testAccount := &models.Account{
		ID:          999,
		UserID:      1,
		AccountType: models.Checking,
		Balance:     100.0,
	}
	
	// Set up expectations for account not found
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("account not found"))
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.UpdateAccount(testAccount)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())
	mockRepo.AssertNotCalled(t, "Update")
	mockRepo.AssertExpectations(t)
}

func TestDeleteAccount_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up test account
	testAccount := &models.Account{
		ID:          1,
		UserID:      1,
		AccountType: models.Checking,
		Balance:     100.0,
	}
	
	// Set up expectations
	mockRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.DeleteAccount(1)
	
	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAccount_NotFound(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockAccountRepository)
	
	// Set up expectations for account not found
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("account not found"))
	
	// Create service with mock repo
	service := NewAccountService(mockRepo)
	
	// Call the method being tested
	err := service.DeleteAccount(999)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())
	mockRepo.AssertNotCalled(t, "Delete")
	mockRepo.AssertExpectations(t)
}
