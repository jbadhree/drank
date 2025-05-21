package repository

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id string) (*models.Account, error)
	FindByUserID(userID string) ([]models.Account, error)
	FindByAccountNumber(accountNumber string) (*models.Account, error)
	FindAll() ([]models.Account, error)
	Update(account *models.Account) error
	UpdateWithTx(account *models.Account) error
	Delete(id string) error
	UpdateBalance(id string, amount float64) error
	FindByIDWithLock(id string) (*models.Account, error)
}

type accountRepository struct {
	client *firestore.Client
	ctx    context.Context
}

func NewAccountRepository(client *firestore.Client) AccountRepository {
	return &accountRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *accountRepository) Create(account *models.Account) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = account.CreatedAt
	docRef, _, err := r.client.Collection("accounts").Add(r.ctx, account)
	if err != nil {
		return err
	}
	account.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
	return err
}

func (r *accountRepository) FindByID(id string) (*models.Account, error) {
	doc, err := r.client.Collection("accounts").Doc(id).Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	var accountData models.Account
	err = doc.DataTo(&accountData)
	if err != nil {
		return nil, err
	}
	return &accountData, nil
}

func (r *accountRepository) FindByUserID(userID string) ([]models.Account, error) {
	var accounts []models.Account
	iter := r.client.Collection("accounts").Where("userId", "==", userID).Documents(r.ctx)

	// iterate over the documents
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, errors.New("account not found")
			}
			return nil, err
		}
		var account models.Account
		err = doc.DataTo(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *accountRepository) FindByAccountNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	iter := r.client.Collection("accounts").Where("accountNumber", "==", accountNumber).Limit(1).Documents(r.ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	err = doc.DataTo(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindAll() ([]models.Account, error) {
	var accounts []models.Account
	iter := r.client.Collection("accounts").Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, errors.New("account not found")
			}
			return nil, err
		}
		var account models.Account
		err = doc.DataTo(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *accountRepository) Update(account *models.Account) error {
	if account.ID == "" {
		return errors.New("account ID is required for update")
	}
	account.UpdatedAt = time.Now()
	_, err := r.client.Collection("accounts").Doc(account.ID).Set(r.ctx, account)
	return err
}

func (r *accountRepository) UpdateWithTx(account *models.Account) error {
	if account.ID == "" {
		return errors.New("account ID is required for update")
	}
	err := r.client.RunTransaction(r.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := r.client.Collection("accounts").Doc(account.ID)
		account.UpdatedAt = time.Now()
		if err := tx.Set(docRef, account); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *accountRepository) Delete(id string) error {
	_, err := r.client.Collection("accounts").Doc(id).Delete(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errors.New("account not found")
		}
		return err
	}
	return nil
}

func (r *accountRepository) UpdateBalance(id string, amount float64) error {
	account, err := r.FindByID(id)
	if err != nil {
		return err
	}
	account.Balance += amount
	account.UpdatedAt = time.Now()
	_, err = r.client.Collection("accounts").Doc(id).Set(r.ctx, account)
	return err
}

// FindByIDWithLock reads an account document inside a Firestore transaction for safe concurrent updates.
func (r *accountRepository) FindByIDWithLock(id string) (*models.Account, error) {

	var account models.Account
	err := r.client.RunTransaction(r.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := r.client.Collection("accounts").Doc(id)
		docSnap, err := tx.Get(docRef)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return errors.New("account not found")
			}
			return err
		}
		err = docSnap.DataTo(&account)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &account, nil
}
