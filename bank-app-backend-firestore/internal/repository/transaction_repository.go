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

// TransactionRepositoryImpl - Implementation of the TransactionRepository interface
type TransactionRepositoryImpl struct {
	client *firestore.Client
	ctx    context.Context
	userID string
}

// NewTransactionRepository - Create a new transaction repository
func NewTransactionRepository(client *firestore.Client, userID string) interfaces.TransactionRepository {
	return &TransactionRepositoryImpl{
		client: client,
		ctx:    context.Background(),
		userID: userID,
	}
}

// getCollectionName returns the user-prefixed collection name
func (r *TransactionRepositoryImpl) getCollectionName() string {
	return r.userID + "_transactions"
}

// Create - Create a new transaction
func (r *TransactionRepositoryImpl) Create(transaction models.Transaction) (models.Transaction, error) {
	// Set created and updated timestamps
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = now
	}

	// Add transaction to Firestore
	docRef, _, err := r.client.Collection(r.getCollectionName()).Add(r.ctx, transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	// Update the transaction with the generated ID
	transaction.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// FindByID - Find transaction by ID
func (r *TransactionRepositoryImpl) FindByID(id string) (models.Transaction, error) {
	docRef := r.client.Collection(r.getCollectionName()).Doc(id)
	docSnapshot, err := docRef.Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.Transaction{}, errors.New("transaction not found")
		}
		return models.Transaction{}, err
	}

	var transaction models.Transaction
	err = docSnapshot.DataTo(&transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// FindByAccountID - Find transactions by account ID
func (r *TransactionRepositoryImpl) FindByAccountID(accountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.client.Collection(r.getCollectionName()).Where("accountId", "==", accountID).OrderBy("transactionDate", firestore.Desc)
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

		var transaction models.Transaction
		err = doc.DataTo(&transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindAll - Find all transactions
func (r *TransactionRepositoryImpl) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction

	iter := r.client.Collection(r.getCollectionName()).OrderBy("transactionDate", firestore.Desc).Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var transaction models.Transaction
		err = doc.DataTo(&transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindBySourceAccountID - Find transactions by source account ID
func (r *TransactionRepositoryImpl) FindBySourceAccountID(sourceAccountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.client.Collection(r.getCollectionName()).Where("sourceAccountId", "==", sourceAccountID).OrderBy("transactionDate", firestore.Desc)
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

		var transaction models.Transaction
		err = doc.DataTo(&transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindByTargetAccountID - Find transactions by target account ID
func (r *TransactionRepositoryImpl) FindByTargetAccountID(targetAccountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.client.Collection(r.getCollectionName()).Where("targetAccountId", "==", targetAccountID).OrderBy("transactionDate", firestore.Desc)
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

		var transaction models.Transaction
		err = doc.DataTo(&transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// CreateTransfer - Create a transfer transaction between two accounts using Firestore transaction
func (r *TransactionRepositoryImpl) CreateTransfer(sourceAccountID, targetAccountID string, amount float64, description string) error {
	sourceAccountRef := r.client.Collection(r.userID + "_accounts").Doc(sourceAccountID)
	targetAccountRef := r.client.Collection(r.userID + "_accounts").Doc(targetAccountID)

	// Use a Firestore transaction for atomic operation
	return r.client.RunTransaction(r.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get source account
		sourceAccountDoc, err := tx.Get(sourceAccountRef)
		if err != nil {
			return err
		}
		var sourceAccount models.Account
		if err := sourceAccountDoc.DataTo(&sourceAccount); err != nil {
			return err
		}

		// Get target account
		targetAccountDoc, err := tx.Get(targetAccountRef)
		if err != nil {
			return err
		}
		var targetAccount models.Account
		if err := targetAccountDoc.DataTo(&targetAccount); err != nil {
			return err
		}

		// Check if source account has enough balance
		if sourceAccount.Balance < amount {
			return errors.New("insufficient balance in source account")
		}

		// Update account balances
		now := time.Now()
		sourceAccount.Balance -= amount
		sourceAccount.UpdatedAt = now
		targetAccount.Balance += amount
		targetAccount.UpdatedAt = now

		// Create source account transaction (withdrawal)
		sourceTransactionRef := r.client.Collection(r.getCollectionName()).NewDoc()
		sourceTransaction := models.Transaction{
			ID:              sourceTransactionRef.ID,
			AccountID:       sourceAccountID,
			SourceAccountID: &sourceAccountID,
			TargetAccountID: &targetAccountID,
			Amount:          -amount,
			Balance:         sourceAccount.Balance,
			Type:            models.Transfer,
			Description:     description,
			TransactionDate: now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		// Create target account transaction (deposit)
		targetTransactionRef := r.client.Collection(r.getCollectionName()).NewDoc()
		targetTransaction := models.Transaction{
			ID:              targetTransactionRef.ID,
			AccountID:       targetAccountID,
			SourceAccountID: &sourceAccountID,
			TargetAccountID: &targetAccountID,
			Amount:          amount,
			Balance:         targetAccount.Balance,
			Type:            models.Transfer,
			Description:     description,
			TransactionDate: now,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		// Update accounts and create transactions in the transaction
		tx.Set(sourceAccountRef, sourceAccount)
		tx.Set(targetAccountRef, targetAccount)
		tx.Set(sourceTransactionRef, sourceTransaction)
		tx.Set(targetTransactionRef, targetTransaction)

		return nil
	})
}
