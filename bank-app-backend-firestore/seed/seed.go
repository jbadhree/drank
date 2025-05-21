package seed

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"google.golang.org/api/iterator"
)

// SeedDatabase seeds the 'drank' database with initial data
func SeedDatabase(client *firestore.Client) error {
	// Clear existing data
	if err := clearData(client); err != nil {
		return err
	}

	// Seed users
	users, err := seedUsers(client)
	if err != nil {
		return err
	}

	// Seed accounts
	accounts, err := seedAccounts(client, users)
	if err != nil {
		return err
	}

	// Seed transactions
	if err := seedTransactions(client, accounts); err != nil {
		return err
	}

	return nil
}

func clearData(client *firestore.Client) error {
	ctx := context.Background()
	collections := []string{"transactions", "accounts", "users"}
	for _, col := range collections {
		iter := client.Collection(col).Documents(ctx)
		batch := client.Batch()
		count := 0
		for {
			doc, err := iter.Next()
			if err != nil {
				if err.Error() == "iterator.Done" || err == iterator.Done {
					break
				}
				return err
			}
			batch.Delete(doc.Ref)
			count++
			// Firestore limits batches to 500 operations
			if count == 500 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return err
				}
				batch = client.Batch()
				count = 0
			}
		}
		if count > 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func seedUsers(client *firestore.Client) ([]models.User, error) {
	users := []models.User{
		{
			Email:     "john.doe@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "jane.smith@example.com",
			Password:  "password123",
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}

	for i := range users {
		// Add to Firestore and get the generated ID
		// Hash the password before storing
		err := users[i].BeforeSave()
		if err != nil {
			return nil, err
		}

		// Add user to Firestore
		docRef, _, err := client.Collection("users").Add(context.Background(), users[i])
		if err != nil {
			return nil, err
		}
		users[i].ID = docRef.ID
		// Optionally update the document with the ID field
		_, err = docRef.Set(context.Background(), map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func seedAccounts(client *firestore.Client, users []models.User) ([]models.Account, error) {
	ctx := context.Background()
	accounts := []models.Account{}

	accountTypes := []models.AccountType{models.Checking, models.Savings}

	// Create checking and savings accounts for each user
	for _, user := range users {
		for _, accountType := range accountTypes {
			account := models.Account{
				UserID:        user.ID,
				AccountNumber: generateAccountNumber(),
				AccountType:   accountType,
				Balance:       5000.00, // Initial balance
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			// Add to Firestore and get the generated ID
			docRef, _, err := client.Collection("accounts").Add(ctx, account)
			if err != nil {
				return nil, err
			}
			account.ID = docRef.ID
			// Optionally update the document with the ID field
			_, err = docRef.Set(ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func seedTransactions(client *firestore.Client, accounts []models.Account) error {

	ctx := context.Background()
	// Transaction types
	transactionTypes := []models.TransactionType{models.Deposit, models.Withdrawal}

	// Description templates
	depositDescriptions := []string{
		"Salary deposit",
		"Refund",
		"Interest earned",
		"Client payment",
		"Tax return",
	}

	withdrawalDescriptions := []string{
		"ATM withdrawal",
		"Online purchase",
		"Bill payment",
		"Subscription payment",
		"Rent payment",
	}

	// Seed random transactions for each account
	for _, account := range accounts {
		// Initial deposit to set up the account
		initialDeposit := models.Transaction{
			AccountID:       account.ID,
			Amount:          account.Balance,
			Balance:         account.Balance,
			Type:            models.Deposit,
			Description:     "Initial deposit",
			TransactionDate: time.Now().AddDate(0, 0, -30), // 30 days ago
		}

		// Create initial deposit transaction
		docRef, _, err := client.Collection("transactions").Add(ctx, initialDeposit)
		if err != nil {
			return err
		}
		initialDeposit.ID = docRef.ID
		// Optionally update the document with the ID field
		_, err = docRef.Set(ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
		if err != nil {
			return err
		}

		// Add 10-15 random transactions for each account
		numTransactions := rand.Intn(6) + 10 // 10-15 transactions

		balance := account.Balance

		for i := 0; i < numTransactions; i++ {
			// Randomize transaction type
			txType := transactionTypes[rand.Intn(len(transactionTypes))]

			// Randomize amount between $10 and $1000
			amount := 10.0 + rand.Float64()*990.0
			amount = float64(int(amount*100)) / 100 // Round to 2 decimal places

			// Handle balance changes based on transaction type
			var description string
			if txType == models.Deposit {
				balance += amount
				description = depositDescriptions[rand.Intn(len(depositDescriptions))]
			} else {
				// For withdrawals, ensure we don't go below zero
				if balance < amount {
					amount = balance * 0.5 // Take only half of what's left
				}
				balance -= amount
				description = withdrawalDescriptions[rand.Intn(len(withdrawalDescriptions))]
			}

			// Create transaction
			transaction := models.Transaction{
				AccountID:       account.ID,
				Amount:          amount,
				Balance:         balance,
				Type:            txType,
				Description:     description,
				TransactionDate: time.Now().AddDate(0, 0, -rand.Intn(30)), // Random date within the last 30 days
			}

			// Create transaction in Firestore
			docRef, _, err := client.Collection("transactions").Add(ctx, transaction)
			if err != nil {
				return err
			}
			transaction.ID = docRef.ID
			// Optionally update the document with the ID field
			_, err = docRef.Set(ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
			if err != nil {
				return err
			}
		}

		// Update the account balance to match final transaction
		_, err = client.Collection("accounts").Doc(account.ID).Set(ctx, map[string]interface{}{"balance": balance}, firestore.MergeAll)
		if err != nil {
			return err
		}

	}

	// Add some transfers between accounts
	if len(accounts) >= 2 {
		// For each user, create a transfer between their checking and savings accounts
		for i := 0; i < len(accounts); i += 2 {
			if i+1 < len(accounts) {
				// Get source and target accounts
				fromAccount := accounts[i]
				toAccount := accounts[i+1]

				// Transfer amount
				amount := 500.0

				// Update balances
				fromBalance := fromAccount.Balance - amount
				toBalance := toAccount.Balance + amount

				// Create withdrawal transaction for source account
				withdrawalTx := models.Transaction{
					AccountID:       fromAccount.ID,
					SourceAccountID: &fromAccount.ID,
					TargetAccountID: &toAccount.ID,
					Amount:          amount,
					Balance:         fromBalance,
					Type:            models.Transfer,
					Description:     fmt.Sprintf("Transfer to account %s", toAccount.AccountNumber),
					TransactionDate: time.Now().AddDate(0, 0, -5), // 5 days ago
				}

				// Create deposit transaction for target account
				depositTx := models.Transaction{
					AccountID:       toAccount.ID,
					SourceAccountID: &fromAccount.ID,
					TargetAccountID: &toAccount.ID,
					Amount:          amount,
					Balance:         toBalance,
					Type:            models.Transfer,
					Description:     fmt.Sprintf("Transfer from account %s", fromAccount.AccountNumber),
					TransactionDate: time.Now().AddDate(0, 0, -5), // 5 days ago
				}

				// Save transactions
				withdrawalDocRef, _, err := client.Collection("transactions").Add(ctx, withdrawalTx)
				if err != nil {
					return err
				}
				withdrawalTx.ID = withdrawalDocRef.ID
				// Optionally update the document with the ID field
				_, err = withdrawalDocRef.Set(ctx, map[string]interface{}{"id": withdrawalDocRef.ID}, firestore.MergeAll)
				if err != nil {
					return err
				}

				// Save deposit transaction
				depositDocRef, _, err := client.Collection("transactions").Add(ctx, depositTx)
				if err != nil {
					return err
				}
				depositTx.ID = depositDocRef.ID
				// Optionally update the document with the ID field
				_, err = depositDocRef.Set(ctx, map[string]interface{}{"id": depositDocRef.ID}, firestore.MergeAll)
				if err != nil {
					return err
				}

				// Update account balances

				_, err = client.Collection("accounts").Doc(fromAccount.ID).Set(ctx, map[string]interface{}{"balance": fromBalance}, firestore.MergeAll)
				if err != nil {
					return err
				}

				_, err = client.Collection("accounts").Doc(toAccount.ID).Set(ctx, map[string]interface{}{"balance": toBalance}, firestore.MergeAll)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func generateAccountNumber() string {
	// Generate a random 10-digit account number
	rand.Seed(time.Now().UnixNano())
	accountNumber := ""
	for i := 0; i < 10; i++ {
		accountNumber += fmt.Sprintf("%d", rand.Intn(10))
	}
	return accountNumber
}
