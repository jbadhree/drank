package models

import (
	"time"
)

type AccountType string

const (
	Checking AccountType = "CHECKING"
	Savings  AccountType = "SAVINGS"
)

// Account - Account model for Firestore
type Account struct {
	ID            string      `json:"id" firestore:"id"`
	UserID        string      `json:"userId" firestore:"userId"`
	AccountNumber string      `json:"accountNumber" firestore:"accountNumber"`
	AccountType   AccountType `json:"accountType" firestore:"accountType"`
	Balance       float64     `json:"balance" firestore:"balance"`
	CreatedAt     time.Time   `json:"createdAt" firestore:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt" firestore:"updatedAt"`
}

// AccountDTO - Data Transfer Object for Account
type AccountDTO struct {
	ID            string      `json:"id"`
	UserID        string      `json:"userId"`
	AccountNumber string      `json:"accountNumber"`
	AccountType   AccountType `json:"accountType"`
	Balance       float64     `json:"balance"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

// ToDTO - Convert Account model to DTO
func (a *Account) ToDTO() AccountDTO {
	return AccountDTO{
		ID:            a.ID,
		UserID:        a.UserID,
		AccountNumber: a.AccountNumber,
		AccountType:   a.AccountType,
		Balance:       a.Balance,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}
