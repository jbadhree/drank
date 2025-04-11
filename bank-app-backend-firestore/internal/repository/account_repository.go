package repository

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountRepositoryImpl - Implementation of the AccountRepository interface
type AccountRepositoryImpl struct {
	client *firestore.Client
	ctx    context.Context
}

// NewAccountRepository - Create a new account repository
func NewAccountRepository(client *firestore.Client) interfaces.AccountRepository {
	return &AccountRepositoryImpl{
		client: client,
		ctx:    context.Background(),
	}
}

// Create - Create a new account
func (r *AccountRepositoryImpl) Create(account models.Account) (models.Account, error) {
	// Check if account number already exists
	query := r.client.Collection("accounts").Where("accountNumber", "==", account.AccountNumber).Limit(1)
	iter := query.Documents(r.ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err == nil {
		return models.Account{}, errors.New("account with this account number already exists")
	}
	if err != iterator.Done {
		return models.Account{}, err
	}

	// Set created and updated timestamps
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	// Add account to Firestore
	docRef, _, err := r.client.Collection("accounts").Add(r.ctx, account)
	if err != nil {
		return models.Account{}, err
	}

	// Update the account with the generated ID
	account.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

// FindByID - Find account by ID
func (r *AccountRepositoryImpl) FindByID(id string) (models.Account, error) {
	docRef := r.client.Collection("accounts").Doc(id)
	docSnapshot, err := docRef.Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}

	var account models.Account
	err = docSnapshot.DataTo(&account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

// FindByAccountNumber - Find account by account number
func (r *AccountRepositoryImpl) FindByAccountNumber(accountNumber string) (models.Account, error) {
	query := r.client.Collection("accounts").Where("accountNumber", "==", accountNumber).Limit(1)
	iter := query.Documents(r.ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}

	var account models.Account
	err = doc.DataTo(&account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

// FindByUserID - Find accounts by user ID
func (r *AccountRepositoryImpl) FindByUserID(userID string) ([]models.Account, error) {
	var accounts []models.Account

	query := r.client.Collection("accounts").Where("userId", "==", userID).OrderBy("createdAt", firestore.Desc)
	iter := query.Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
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

// FindAll - Find all accounts
func (r *AccountRepositoryImpl) FindAll() ([]models.Account, error) {
	var accounts []models.Account

	iter := r.client.Collection("accounts").Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
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

// Update - Update an account
func (r *AccountRepositoryImpl) Update(account models.Account) (models.Account, error) {
	// Check if account exists
	docRef := r.client.Collection("accounts").Doc(account.ID)
	_, err := docRef.Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.Account{}, errors.New("account not found")
		}
		return models.Account{}, err
	}

	// Update account
	account.UpdatedAt = time.Now()
	_, err = docRef.Set(r.ctx, account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

// Delete - Delete an account
func (r *AccountRepositoryImpl) Delete(id string) error {
	_, err := r.client.Collection("accounts").Doc(id).Delete(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBalance - Update account balance
func (r *AccountRepositoryImpl) UpdateBalance(id string, amount float64) (models.Account, error) {
	// Get account
	account, err := r.FindByID(id)
	if err != nil {
		return models.Account{}, err
	}

	// Update balance
	account.Balance += amount
	account.UpdatedAt = time.Now()

	// Update in database
	docRef := r.client.Collection("accounts").Doc(id)
	_, err = docRef.Set(r.ctx, account)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}
