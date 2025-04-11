package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
)

// AccountService - Service for account operations
type AccountService struct {
	repo interfaces.AccountRepository
}

// NewAccountService - Create a new account service
func NewAccountService(repo interfaces.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// Create - Create a new account
func (s *AccountService) Create(account models.Account) (models.AccountDTO, error) {
	// Generate account number if not provided
	if account.AccountNumber == "" {
		account.AccountNumber = generateAccountNumber()
	}

	// Create the account
	createdAccount, err := s.repo.Create(account)
	if err != nil {
		return models.AccountDTO{}, err
	}

	return createdAccount.ToDTO(), nil
}

// GetByID - Get account by ID
func (s *AccountService) GetByID(id string) (models.AccountDTO, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return models.AccountDTO{}, err
	}

	return account.ToDTO(), nil
}

// GetByAccountNumber - Get account by account number
func (s *AccountService) GetByAccountNumber(accountNumber string) (models.AccountDTO, error) {
	account, err := s.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return models.AccountDTO{}, err
	}

	return account.ToDTO(), nil
}

// GetByUserID - Get accounts by user ID
func (s *AccountService) GetByUserID(userID string) ([]models.AccountDTO, error) {
	accounts, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	accountDTOs := make([]models.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}

	return accountDTOs, nil
}

// GetAll - Get all accounts
func (s *AccountService) GetAll() ([]models.AccountDTO, error) {
	accounts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	accountDTOs := make([]models.AccountDTO, len(accounts))
	for i, account := range accounts {
		accountDTOs[i] = account.ToDTO()
	}

	return accountDTOs, nil
}

// Update - Update an account
func (s *AccountService) Update(account models.Account) (models.AccountDTO, error) {
	updatedAccount, err := s.repo.Update(account)
	if err != nil {
		return models.AccountDTO{}, err
	}

	return updatedAccount.ToDTO(), nil
}

// Delete - Delete an account
func (s *AccountService) Delete(id string) error {
	return s.repo.Delete(id)
}

// UpdateBalance - Update account balance
func (s *AccountService) UpdateBalance(id string, amount float64) (models.AccountDTO, error) {
	account, err := s.repo.UpdateBalance(id, amount)
	if err != nil {
		return models.AccountDTO{}, err
	}

	return account.ToDTO(), nil
}

// Helper function to generate a random account number
func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	accountNumber := fmt.Sprintf("%010d", rand.Intn(10000000000))
	return accountNumber
}
