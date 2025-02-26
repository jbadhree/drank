package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/repository"
)

type TransactionService interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error)
	GetAllTransactions(limit, offset int) ([]models.Transaction, error)
	Transfer(request *models.TransferRequest) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, accountRepo repository.AccountRepository) TransactionService {
	return &transactionService{transactionRepo, accountRepo}
}

func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
	// Set transaction date if not provided
	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	// Validate transaction type
	if transaction.Type != models.Deposit && transaction.Type != models.Withdrawal && transaction.Type != models.Transfer {
		return errors.New("invalid transaction type")
	}

	// Get account to update balance
	account, err := s.accountRepo.FindByID(transaction.AccountID)
	if err != nil {
		return err
	}

	// Update account balance based on transaction type
	switch transaction.Type {
	case models.Deposit:
		account.Balance += transaction.Amount
	case models.Withdrawal:
		if account.Balance < transaction.Amount {
			return errors.New("insufficient funds")
		}
		account.Balance -= transaction.Amount
	case models.Transfer:
		// Transfers are handled in the Transfer method
		return errors.New("use transfer method for transfer transactions")
	}

	// Set the balance after transaction
	transaction.Balance = account.Balance

	// Update account balance
	if err := s.accountRepo.Update(account); err != nil {
		return err
	}

	// Create transaction
	if err := s.transactionRepo.Create(transaction); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id)
}

func (s *transactionService) GetTransactionsByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error) {
	return s.transactionRepo.FindByAccountID(accountID, limit, offset)
}

func (s *transactionService) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	return s.transactionRepo.FindAll(limit, offset)
}

func (s *transactionService) Transfer(request *models.TransferRequest) error {
	if request.Amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	if request.FromAccountID == request.ToAccountID {
		return errors.New("cannot transfer to the same account")
	}

	// Lock the from account for update to prevent race conditions
	fromAccount, fromTx, err := s.accountRepo.FindByIDWithLock(request.FromAccountID)
	if err != nil {
		return err
	}
	defer func() {
		if fromTx != nil {
			fromTx.Rollback()
		}
	}()

	// Check if from account has sufficient balance
	if fromAccount.Balance < request.Amount {
		return errors.New("insufficient funds")
	}

	// Lock the to account for update
	toAccount, toTx, err := s.accountRepo.FindByIDWithLock(request.ToAccountID)
	if err != nil {
		return err
	}
	defer func() {
		if toTx != nil {
			toTx.Rollback()
		}
	}()

	// Update balances
	fromAccount.Balance -= request.Amount
	toAccount.Balance += request.Amount

	// Save from account
	result := fromTx.Save(fromAccount)
	if err := result.Error(); err != nil {
		return err
	}

	// Save to account
	result = toTx.Save(toAccount)
	if err := result.Error(); err != nil {
		return err
	}

	// Create withdrawal transaction for from account
	withdrawalDesc := fmt.Sprintf("Transfer to account %s: %s", toAccount.AccountNumber, request.Description)
	withdrawal := &models.Transaction{
		AccountID:       fromAccount.ID,
		SourceAccountID: &fromAccount.ID,
		TargetAccountID: &toAccount.ID,
		Amount:          request.Amount,
		Balance:         fromAccount.Balance,
		Type:            models.Transfer,
		Description:     withdrawalDesc,
		TransactionDate: time.Now(),
	}

	// Create deposit transaction for to account
	depositDesc := fmt.Sprintf("Transfer from account %s: %s", fromAccount.AccountNumber, request.Description)
	deposit := &models.Transaction{
		AccountID:       toAccount.ID,
		SourceAccountID: &fromAccount.ID,
		TargetAccountID: &toAccount.ID,
		Amount:          request.Amount,
		Balance:         toAccount.Balance,
		Type:            models.Transfer,
		Description:     depositDesc,
		TransactionDate: time.Now(),
	}

	// Create withdrawal transaction
	if err := s.transactionRepo.CreateWithTx(withdrawal, fromTx); err != nil {
		return err
	}

	// Create deposit transaction
	if err := s.transactionRepo.CreateWithTx(deposit, toTx); err != nil {
		return err
	}

	// Commit from account transaction
	commitResult := fromTx.Commit()
	if err := commitResult.Error(); err != nil {
		return err
	}
	fromTx = nil

	// Commit to account transaction
	commitResult = toTx.Commit()
	if err := commitResult.Error(); err != nil {
		return err
	}
	toTx = nil

	return nil
}
