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
func SeedDatabase(client *firestore.Client, userID string) error {
	ctx := context.Background()

	usersCol := userID + "_users"
	accountsCol := userID + "_accounts"
	transactionsCol := userID + "_transactions"

	// Clear existing data in collections before seeding using BulkWriter
	collections := []string{usersCol, accountsCol, transactionsCol}
	for _, col := range collections {
		colRef := client.Collection(col)
		for {
			docs, err := colRef.Limit(100).Documents(ctx).GetAll()
			if err != nil {
				return err
			}
			if len(docs) == 0 {
				break
			}

			bulkWriter := client.BulkWriter(ctx)
			for _, doc := range docs {
				bulkWriter.Delete(doc.Ref)
			}
			bulkWriter.End()
		}
	}

	// Check if users collection already has data
	usersRef := client.Collection(usersCol)
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
			Email:     "john.doe@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:     "jane.smith@example.com",
			Password:  "password123",
			FirstName: "Jane",
			LastName:  "Smith",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Encrypt passwords before inserting users
	for i, user := range users {
		hashed, err := models.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}
		users[i].Password = hashed
	}

	// Create users in batch using BulkWriter
	bulkWriter := client.BulkWriter(ctx)
	userIDs := make([]string, len(users))

	for i, user := range users {
		userRef := client.Collection(usersCol).NewDoc()
		userIDs[i] = userRef.ID
		user.ID = userRef.ID
		bulkWriter.Set(userRef, user)
	}
	bulkWriter.End()

	// Create accounts for the users using BulkWriter
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
	bulkWriter = client.BulkWriter(ctx)
	accountIDs := make([]string, len(accounts))

	for i, account := range accounts {
		accountRef := client.Collection(accountsCol).NewDoc()
		accountIDs[i] = accountRef.ID
		account.ID = accountRef.ID
		bulkWriter.Set(accountRef, account)
	}
	bulkWriter.End()

	// Create some sample transactions using BulkWriter
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
	bulkWriter = client.BulkWriter(ctx)

	for _, transaction := range transactions {
		transactionRef := client.Collection(transactionsCol).NewDoc()
		transaction.ID = transactionRef.ID
		bulkWriter.Set(transactionRef, transaction)
	}
	bulkWriter.End()

	log.Println("Database seeded successfully!")
	return nil
}
