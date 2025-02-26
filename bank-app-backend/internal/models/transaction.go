package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
)

type Transaction struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	AccountID        uint             `json:"accountId" gorm:"not null"`
	SourceAccountID  *uint            `json:"sourceAccountId,omitempty"`
	TargetAccountID  *uint            `json:"targetAccountId,omitempty"`
	Amount           float64          `json:"amount" gorm:"not null"`
	Balance          float64          `json:"balance" gorm:"not null"` // Balance after the transaction
	Type             TransactionType  `json:"type" gorm:"not null"`
	Description      string           `json:"description"`
	TransactionDate  time.Time        `json:"transactionDate" gorm:"not null"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt   `json:"-" gorm:"index"`
}

// TransactionDTO - Data Transfer Object for Transaction
type TransactionDTO struct {
	ID               uint            `json:"id"`
	AccountID        uint            `json:"accountId"`
	SourceAccountID  *uint           `json:"sourceAccountId,omitempty"`
	TargetAccountID  *uint           `json:"targetAccountId,omitempty"`
	Amount           float64         `json:"amount"`
	Balance          float64         `json:"balance"`
	Type             TransactionType `json:"type"`
	Description      string          `json:"description"`
	TransactionDate  time.Time       `json:"transactionDate"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

// ToDTO - Convert Transaction model to DTO
func (t *Transaction) ToDTO() TransactionDTO {
	return TransactionDTO{
		ID:               t.ID,
		AccountID:        t.AccountID,
		SourceAccountID:  t.SourceAccountID,
		TargetAccountID:  t.TargetAccountID,
		Amount:           t.Amount,
		Balance:          t.Balance,
		Type:             t.Type,
		Description:      t.Description,
		TransactionDate:  t.TransactionDate,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

// TransferRequest - Request body for transfer
type TransferRequest struct {
	FromAccountID uint    `json:"fromAccountId" binding:"required"`
	ToAccountID   uint    `json:"toAccountId" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
}
