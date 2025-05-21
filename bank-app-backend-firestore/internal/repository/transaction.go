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

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id string) (*models.Transaction, error)
	FindByAccountID(accountID string, limit, offset int) ([]models.Transaction, error)
	FindAll(limit, offset int) ([]models.Transaction, error)
	CountByAccountID(accountID string) (int64, error)
	CountAll() (int64, error)
	CreateWithTx(transaction *models.Transaction) error
}

type transactionRepository struct {
	client *firestore.Client
	ctx    context.Context
}

func NewTransactionRepository(client *firestore.Client) TransactionRepository {
	return &transactionRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = transaction.CreatedAt
	docRef, _, err := r.client.Collection("transactions").Add(r.ctx, transaction)
	if err != nil {
		return err
	}
	transaction.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
	return err
}

func (r *transactionRepository) FindByID(id string) (*models.Transaction, error) {
	doc, err := r.client.Collection("transactions").Doc(id).Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	var transactionData models.Transaction
	err = doc.DataTo(&transactionData)
	if err != nil {
		return nil, err
	}
	return &transactionData, nil
}

func (r *transactionRepository) FindByAccountID(accountID string, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	iter := r.client.Collection("transactions").Where("accountId", "==", accountID).OrderBy("transactionDate", firestore.Desc).Limit(limit).Offset(offset).Documents(r.ctx)

	// iterate over the documents
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, errors.New("transaction not found")
			}
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

func (r *transactionRepository) FindAll(limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	iter := r.client.Collection("transactions").Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, errors.New("transaction not found")
			}
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

func (r *transactionRepository) CountByAccountID(accountID string) (int64, error) {
	var count int64
	iter := r.client.Collection("transactions").Where("accountId", "==", accountID).Documents(r.ctx)
	defer iter.Stop()

	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

func (r *transactionRepository) CountAll() (int64, error) {
	var count int64
	iter := r.client.Collection("transactions").Documents(r.ctx)
	defer iter.Stop()

	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

// CreateWithTx creates a transaction document within a Firestore transaction.
func (r *transactionRepository) CreateWithTx(transaction *models.Transaction) error {
	err := r.client.RunTransaction(r.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := r.client.Collection("transactions").NewDoc()
		transaction.ID = docRef.ID
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = transaction.CreatedAt
		if err := tx.Set(docRef, transaction); err != nil {
			return err
		}
		return nil
	})
	return err
}
