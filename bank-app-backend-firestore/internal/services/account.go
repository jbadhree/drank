package services

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/repository"
)

type AccountService interface {
	CreateAccount(account *models.Account) error
	GetAccountByID(id string) (*models.Account, error)
	GetAccountsByUserID(userID string) ([]models.Account, error)
	GetAllAccounts() ([]models.Account, error)
	UpdateAccount(account *models.Account) error
	DeleteAccount(id string) error
	GenerateAccountNumber() string
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo}
}

func (s *accountService) CreateAccount(account *models.Account) error {
	// Generate account number if not provided
	if account.AccountNumber == "" {
		account.AccountNumber = s.GenerateAccountNumber()
	}

	// Validate account type
	if account.AccountType != models.Checking && account.AccountType != models.Savings {
		return errors.New("invalid account type")
	}

	return s.accountRepo.Create(account)
}

func (s *accountService) GetAccountByID(id string) (*models.Account, error) {
	return s.accountRepo.FindByID(id)
}

func (s *accountService) GetAccountsByUserID(userID string) ([]models.Account, error) {
	return s.accountRepo.FindByUserID(userID)
}

func (s *accountService) GetAllAccounts() ([]models.Account, error) {
	return s.accountRepo.FindAll()
}

func (s *accountService) UpdateAccount(account *models.Account) error {
	// Check if the account exists
	_, err := s.accountRepo.FindByID(account.ID)
	if err != nil {
		return err
	}

	return s.accountRepo.Update(account)
}

func (s *accountService) DeleteAccount(id string) error {
	// Check if the account exists
	_, err := s.accountRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.accountRepo.Delete(id)
}

func (s *accountService) GenerateAccountNumber() string {
	// Generate a random 10-digit account number
	rand.Seed(time.Now().UnixNano())
	accountNumber := ""
	for i := 0; i < 10; i++ {
		accountNumber += strconv.Itoa(rand.Intn(10))
	}
	return accountNumber
}
