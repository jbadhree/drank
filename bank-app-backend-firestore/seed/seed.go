package seed

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SeedDatabase - Seed the database with initial data
func SeedDatabase(client *firestore.Client) error {
	ctx := context.Background()

	// Check if users collection already has data
	usersRef := client.Collection("users")
	usersDocs, err := usersRef.Limit(1).Documents(ctx).GetAll()
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}
	
	// If users collection is not empty, don't seed
	if len(usersDocs) > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database...")

	// Create demo users
	users := []models.User{
		{
			Email:     "user1@example.com",
			Password:  "$2a$10$1JGBZUOPMvmbvnAZWF1mTuc0r3lqMANdLl0qQ2pZj5.vcnUo4q0qK", // Password: password1
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:     "user2@example.com",
			Password:  "$2a$10$1JGBZUOPMvmbvnAZWF1mTuc0r3lqMANdLl0qQ2pZj5.vcnUo4q0qK", // Password: password1
			FirstName: "Jane",
			LastName:  "Smith",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Create users in batch
	batch := client.Batch()
	userIDs := make([]string, len(users))

	for i, user := range users {
		userRef := client.Collection("users").NewDoc()
		userIDs[i] = userRef.ID
		user.ID = userRef.ID
		batch.Set(userRef, user)
	}

	// Execute the batch
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}

	// Create accounts for the users
	accounts := []models.Account{
		{
			UserID:        userIDs[0],
			AccountNumber: "1000000001",
			AccountType:   models.Checking,
			Balance:       1000.00,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			UserID:        userIDs[0],
			AccountNumber: "1000000002",
			AccountType:   models.Savings,
			Balance:       5000.00,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			UserID:        userIDs[1],
			AccountNumber: "1000000003",
			AccountType:   models.Checking,
			Balance:       2000.00,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Create accounts in batch
	batch = client.Batch()
	accountIDs := make([]string, len(accounts))

	for i, account := range accounts {
		accountRef := client.Collection("accounts").NewDoc()
		accountIDs[i] = accountRef.ID
		account.ID = accountRef.ID
		batch.Set(accountRef, account)
	}

	// Execute the batch
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}

	// Create some sample transactions
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	transactions := []models.Transaction{
		{
			AccountID:       accountIDs[0],
			Amount:          500.00,
			Balance:         500.00,
			Type:            models.Deposit,
			Description:     "Initial deposit",
			TransactionDate: lastWeek,
			CreatedAt:       lastWeek,
			UpdatedAt:       lastWeek,
		},
		{
			AccountID:       accountIDs[0],
			Amount:          -100.00,
			Balance:         400.00,
			Type:            models.Withdrawal,
			Description:     "ATM withdrawal",
			TransactionDate: yesterday,
			CreatedAt:       yesterday,
			UpdatedAt:       yesterday,
		},
		{
			AccountID:       accountIDs[0],
			Amount:          600.00,
			Balance:         1000.00,
			Type:            models.Deposit,
			Description:     "Salary deposit",
			TransactionDate: now,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			AccountID:       accountIDs[1],
			Amount:          5000.00,
			Balance:         5000.00,
			Type:            models.Deposit,
			Description:     "Initial deposit",
			TransactionDate: lastWeek,
			CreatedAt:       lastWeek,
			UpdatedAt:       lastWeek,
		},
		{
			AccountID:       accountIDs[2],
			Amount:          2000.00,
			Balance:         2000.00,
			Type:            models.Deposit,
			Description:     "Initial deposit",
			TransactionDate: yesterday,
			CreatedAt:       yesterday,
			UpdatedAt:       yesterday,
		},
	}

	// Create transactions in batch
	batch = client.Batch()

	for _, transaction := range transactions {
		transactionRef := client.Collection("transactions").NewDoc()
		transaction.ID = transactionRef.ID
		batch.Set(transactionRef, transaction)
	}

	// Execute the batch
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}

	log.Println("Database seeded successfully!")
	return nil
}
