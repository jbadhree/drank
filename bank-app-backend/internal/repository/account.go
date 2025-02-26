package repository

import (
	"errors"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id uint) (*models.Account, error)
	FindByUserID(userID uint) ([]models.Account, error)
	FindByAccountNumber(accountNumber string) (*models.Account, error)
	FindAll() ([]models.Account, error)
	Update(account *models.Account) error
	Delete(id uint) error
	UpdateBalance(id uint, amount float64) error
	FindByIDWithLock(id uint) (*models.Account, *gorm.DB, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	result := r.db.First(&account, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) FindByUserID(userID uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) FindByAccountNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	result := r.db.Where("account_number = ?", accountNumber).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) FindAll() ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}

func (r *accountRepository) UpdateBalance(id uint, amount float64) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *accountRepository) FindByIDWithLock(id uint) (*models.Account, *gorm.DB, error) {
	var account models.Account
	tx := r.db.Begin()
	result := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, id)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("account not found")
		}
		return nil, nil, result.Error
	}
	return &account, tx, nil
}
