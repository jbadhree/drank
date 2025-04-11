package models

import (
	"time"
)

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
)

// Transaction - Transaction model for Firestore
type Transaction struct {
	ID               string          `json:"id" firestore:"id"`
	AccountID        string          `json:"accountId" firestore:"accountId"`
	SourceAccountID  *string         `json:"sourceAccountId,omitempty" firestore:"sourceAccountId,omitempty"`
	TargetAccountID  *string         `json:"targetAccountId,omitempty" firestore:"targetAccountId,omitempty"`
	Amount           float64         `json:"amount" firestore:"amount"`
	Balance          float64         `json:"balance" firestore:"balance"` // Balance after the transaction
	Type             TransactionType `json:"type" firestore:"type"`
	Description      string          `json:"description" firestore:"description"`
	TransactionDate  time.Time       `json:"transactionDate" firestore:"transactionDate"`
	CreatedAt        time.Time       `json:"createdAt" firestore:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt" firestore:"updatedAt"`
}

// TransactionDTO - Data Transfer Object for Transaction
type TransactionDTO struct {
	ID               string          `json:"id"`
	AccountID        string          `json:"accountId"`
	SourceAccountID  *string         `json:"sourceAccountId,omitempty"`
	TargetAccountID  *string         `json:"targetAccountId,omitempty"`
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
	FromAccountID string  `json:"fromAccountId" binding:"required"`
	ToAccountID   string  `json:"toAccountId" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
}
