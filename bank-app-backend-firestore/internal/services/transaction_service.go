package services

import (
	"errors"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
)

// TransactionService - Service for transaction operations
type TransactionService struct {
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
}

// NewTransactionService - Create a new transaction service
func NewTransactionService(transactionRepo interfaces.TransactionRepository, accountRepo interfaces.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// Create - Create a new transaction
func (s *TransactionService) Create(transaction models.Transaction) (models.TransactionDTO, error) {
	// Get the account
	account, err := s.accountRepo.FindByID(transaction.AccountID)
	if err != nil {
		return models.TransactionDTO{}, err
	}

	// Update account balance
	newBalance := account.Balance + transaction.Amount
	account.Balance = newBalance
	account.UpdatedAt = time.Now()

	// Update transaction with new balance
	transaction.Balance = newBalance

	// Save transaction
	createdTransaction, err := s.transactionRepo.Create(transaction)
	if err != nil {
		return models.TransactionDTO{}, err
	}

	// Update account in database
	_, err = s.accountRepo.Update(account)
	if err != nil {
		// Here we should ideally roll back the transaction, but for simplicity we'll just return an error
		return models.TransactionDTO{}, errors.New("failed to update account balance: " + err.Error())
	}

	return createdTransaction.ToDTO(), nil
}

// GetByID - Get transaction by ID
func (s *TransactionService) GetByID(id string) (models.TransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return models.TransactionDTO{}, err
	}

	return transaction.ToDTO(), nil
}

// GetByAccountID - Get transactions by account ID
func (s *TransactionService) GetByAccountID(accountID string) ([]models.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	transactionDTOs := make([]models.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		transactionDTOs[i] = transaction.ToDTO()
	}

	return transactionDTOs, nil
}

// GetAll - Get all transactions
func (s *TransactionService) GetAll() ([]models.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	transactionDTOs := make([]models.TransactionDTO, len(transactions))
	for i, transaction := range transactions {
		transactionDTOs[i] = transaction.ToDTO()
	}

	return transactionDTOs, nil
}

// Transfer - Transfer funds between accounts
func (s *TransactionService) Transfer(req models.TransferRequest) error {
	// Validate accounts
	if req.FromAccountID == req.ToAccountID {
		return errors.New("cannot transfer to the same account")
	}

	// Validate amount
	if req.Amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	// Check if accounts exist
	_, err := s.accountRepo.FindByID(req.FromAccountID)
	if err != nil {
		return errors.New("source account not found")
	}

	_, err = s.accountRepo.FindByID(req.ToAccountID)
	if err != nil {
		return errors.New("target account not found")
	}

	// Perform the transfer using transaction repository's atomic transaction function
	return s.transactionRepo.CreateTransfer(req.FromAccountID, req.ToAccountID, req.Amount, req.Description)
}
