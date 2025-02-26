package repository

import (
	"errors"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error)
	FindAll(limit, offset int) ([]models.Transaction, error)
	CountByAccountID(accountID uint) (int64, error)
	CountAll() (int64, error)
	CreateWithTx(transaction *models.Transaction, tx *gorm.DB) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, result.Error
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByAccountID(accountID uint, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Where("account_id = ?", accountID).Order("transaction_date DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	
	return transactions, nil
}

func (r *transactionRepository) FindAll(limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Order("transaction_date DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	
	return transactions, nil
}

func (r *transactionRepository) CountByAccountID(accountID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Transaction{}).Where("account_id = ?", accountID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *transactionRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *transactionRepository) CreateWithTx(transaction *models.Transaction, tx *gorm.DB) error {
	return tx.Create(transaction).Error
}
