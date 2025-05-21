package repository

import (
	"gorm.io/gorm"
)

// GormResult defines an interface with an Error method
type GormResult interface {
	Error() error
}

// GormTx defines common transaction methods from gorm.DB
type GormTx interface {
	Save(value interface{}) GormResult
	Commit() GormResult
	Rollback() GormResult
	Create(value interface{}) GormResult
}

// GormDBWrapper wraps *gorm.DB to implement GormTx
type GormDBWrapper struct {
	DB *gorm.DB
}

// GormResultWrapper wraps *gorm.DB to implement GormResult
type GormResultWrapper struct {
	DB *gorm.DB
}

func (w GormResultWrapper) Error() error {
	return w.DB.Error
}

func (w GormDBWrapper) Save(value interface{}) GormResult {
	return GormResultWrapper{w.DB.Save(value)}
}

func (w GormDBWrapper) Create(value interface{}) GormResult {
	return GormResultWrapper{w.DB.Create(value)}
}

func (w GormDBWrapper) Commit() GormResult {
	return GormResultWrapper{w.DB.Commit()}
}

func (w GormDBWrapper) Rollback() GormResult {
	return GormResultWrapper{w.DB.Rollback()}
}
