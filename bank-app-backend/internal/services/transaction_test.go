package services

import (
	"errors"
	"testing"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Create a mock for the transaction repository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	args := m.Called(transaction)
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

func (m *MockTransactionRepository) CreateWithTx(transaction *models.Transaction, tx *gorm.DB) error {
	args := m.Called(transaction, tx)
	return args.Error(0)
}

func TestCreateTransaction_Deposit(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test data
	testAccount := &models.Account{
		ID:           1,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	depositAmount := 50.0
	
	testTransaction := &models.Transaction{
		AccountID:       1,
		Amount:          depositAmount,
		Type:            models.Deposit,
		TransactionDate: time.Now(),
		Description:     "Test deposit",
	}
	
	// Set up expectations
	mockAccountRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	mockAccountRepo.On("Update", mock.MatchedBy(func(account *models.Account) bool {
		return account.ID == 1 && account.Balance == 150.0 // Original balance + deposit
	})).Return(nil)
	mockTransactionRepo.On("Create", mock.MatchedBy(func(transaction *models.Transaction) bool {
		return transaction.AccountID == 1 && 
			   transaction.Type == models.Deposit && 
			   transaction.Amount == depositAmount &&
			   transaction.Balance == 150.0 // New balance after deposit
	})).Return(nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.CreateTransaction(testTransaction)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 150.0, testTransaction.Balance)
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestCreateTransaction_Withdrawal_Success(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test data
	testAccount := &models.Account{
		ID:           1,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	withdrawalAmount := 50.0
	
	testTransaction := &models.Transaction{
		AccountID:       1,
		Amount:          withdrawalAmount,
		Type:            models.Withdrawal,
		TransactionDate: time.Now(),
		Description:     "Test withdrawal",
	}
	
	// Set up expectations
	mockAccountRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	mockAccountRepo.On("Update", mock.MatchedBy(func(account *models.Account) bool {
		return account.ID == 1 && account.Balance == 50.0 // Original balance - withdrawal
	})).Return(nil)
	mockTransactionRepo.On("Create", mock.MatchedBy(func(transaction *models.Transaction) bool {
		return transaction.AccountID == 1 && 
			   transaction.Type == models.Withdrawal && 
			   transaction.Amount == withdrawalAmount &&
			   transaction.Balance == 50.0 // New balance after withdrawal
	})).Return(nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.CreateTransaction(testTransaction)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 50.0, testTransaction.Balance)
	mockAccountRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestCreateTransaction_Withdrawal_InsufficientFunds(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test data
	testAccount := &models.Account{
		ID:           1,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	withdrawalAmount := 150.0 // More than balance
	
	testTransaction := &models.Transaction{
		AccountID:       1,
		Amount:          withdrawalAmount,
		Type:            models.Withdrawal,
		TransactionDate: time.Now(),
		Description:     "Test withdrawal",
	}
	
	// Set up expectations
	mockAccountRepo.On("FindByID", uint(1)).Return(testAccount, nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.CreateTransaction(testTransaction)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
	mockAccountRepo.AssertExpectations(t)
	// These methods should not be called
	mockAccountRepo.AssertNotCalled(t, "Update")
	mockTransactionRepo.AssertNotCalled(t, "Create")
}

func TestCreateTransaction_InvalidType(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test transaction with invalid type
	testTransaction := &models.Transaction{
		AccountID:       1,
		Amount:          50.0,
		Type:            "INVALID_TYPE",
		TransactionDate: time.Now(),
		Description:     "Test invalid type",
	}
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.CreateTransaction(testTransaction)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "invalid transaction type", err.Error())
	// These methods should not be called
	mockAccountRepo.AssertNotCalled(t, "FindByID")
	mockAccountRepo.AssertNotCalled(t, "Update")
	mockTransactionRepo.AssertNotCalled(t, "Create")
}

func TestGetTransactionByID_Success(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test transaction
	testTransaction := &models.Transaction{
		ID:              1,
		AccountID:       1,
		Amount:          50.0,
		Balance:         150.0,
		Type:            models.Deposit,
		TransactionDate: time.Now(),
		Description:     "Test transaction",
	}
	
	// Set up expectations
	mockTransactionRepo.On("FindByID", uint(1)).Return(testTransaction, nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	transaction, err := service.GetTransactionByID(1)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, testTransaction.ID, transaction.ID)
	assert.Equal(t, testTransaction.Amount, transaction.Amount)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Set up expectations for transaction not found
	mockTransactionRepo.On("FindByID", uint(999)).Return(nil, errors.New("transaction not found"))
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	transaction, err := service.GetTransactionByID(999)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "transaction not found", err.Error())
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetTransactionsByAccountID_Success(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test transactions
	testTransactions := []models.Transaction{
		{ID: 1, AccountID: 1, Amount: 50.0, Type: models.Deposit},
		{ID: 2, AccountID: 1, Amount: 25.0, Type: models.Withdrawal},
	}
	
	// Set up expectations
	mockTransactionRepo.On("FindByAccountID", uint(1), 10, 0).Return(testTransactions, nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	transactions, err := service.GetTransactionsByAccountID(1, 10, 0)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, testTransactions[0].ID, transactions[0].ID)
	assert.Equal(t, testTransactions[1].Type, transactions[1].Type)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetAllTransactions_Success(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create test transactions
	testTransactions := []models.Transaction{
		{ID: 1, AccountID: 1, Amount: 50.0, Type: models.Deposit},
		{ID: 2, AccountID: 2, Amount: 25.0, Type: models.Withdrawal},
	}
	
	// Set up expectations
	mockTransactionRepo.On("FindAll", 10, 0).Return(testTransactions, nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	transactions, err := service.GetAllTransactions(10, 0)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, testTransactions[0].ID, transactions[0].ID)
	assert.Equal(t, testTransactions[1].Type, transactions[1].Type)
	mockTransactionRepo.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Mock transaction
	mockTx := &gorm.DB{}
	
	// Create test accounts
	fromAccount := &models.Account{
		ID:           1,
		AccountNumber: "1234567890",
		Balance:      100.0,
	}
	
	toAccount := &models.Account{
		ID:           2,
		AccountNumber: "0987654321",
		Balance:      50.0,
	}
	
	// Create transfer request
	transferRequest := &models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        25.0,
		Description:   "Test transfer",
	}
	
	// Set up expectations
	mockAccountRepo.On("FindByIDWithLock", uint(1)).Return(fromAccount, mockTx, nil)
	mockAccountRepo.On("FindByIDWithLock", uint(2)).Return(toAccount, mockTx, nil)
	mockTransactionRepo.On("CreateWithTx", mock.AnythingOfType("*models.Transaction"), mockTx).Return(nil).Twice()
	mockTx.On("Save", mock.Anything).Return(mockTx)
	mockTx.On("Commit").Return(mockTx)
	mockTx.On("Error").Return(nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.Transfer(transferRequest)
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 75.0, fromAccount.Balance)  // 100 - 25
	assert.Equal(t, 75.0, toAccount.Balance)    // 50 + 25
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Mock transaction
	mockTx := &gorm.DB{}
	
	// Create test accounts
	fromAccount := &models.Account{
		ID:           1,
		AccountNumber: "1234567890",
		Balance:      20.0,  // Less than transfer amount
	}
	
	// Create transfer request
	transferRequest := &models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        25.0,
		Description:   "Test transfer",
	}
	
	// Set up expectations
	mockAccountRepo.On("FindByIDWithLock", uint(1)).Return(fromAccount, mockTx, nil)
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.Transfer(transferRequest)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}

func TestTransfer_SameAccount(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create transfer request with same account for source and destination
	transferRequest := &models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   1,
		Amount:        25.0,
		Description:   "Test transfer",
	}
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.Transfer(transferRequest)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "cannot transfer to the same account", err.Error())
}

func TestTransfer_InvalidAmount(t *testing.T) {
	// Create mock repositories
	mockTransactionRepo := new(MockTransactionRepository)
	mockAccountRepo := new(MockAccountRepository)
	
	// Create transfer request with negative amount
	transferRequest := &models.TransferRequest{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        -25.0,
		Description:   "Test transfer",
	}
	
	// Create service with mock repos
	service := NewTransactionService(mockTransactionRepo, mockAccountRepo)
	
	// Call the method being tested
	err := service.Transfer(transferRequest)
	
	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, "transfer amount must be positive", err.Error())
}
