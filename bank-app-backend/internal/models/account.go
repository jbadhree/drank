package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType string

const (
	Checking AccountType = "CHECKING"
	Savings  AccountType = "SAVINGS"
)

type Account struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"userId" gorm:"not null"`
	AccountNumber string         `json:"accountNumber" gorm:"uniqueIndex;not null"`
	AccountType   AccountType    `json:"accountType" gorm:"not null"`
	Balance       float64        `json:"balance" gorm:"not null;default:0"`
	Transactions  []Transaction  `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// AccountDTO - Data Transfer Object for Account
type AccountDTO struct {
	ID            uint        `json:"id"`
	UserID        uint        `json:"userId"`
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
