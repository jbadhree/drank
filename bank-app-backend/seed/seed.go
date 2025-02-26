package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"gorm.io/gorm"
)

// SeedDatabase seeds the 'drank' database with initial data
func SeedDatabase(db *gorm.DB) error {
	// Clear existing data
	if err := clearData(db); err != nil {
		return err
	}

	// Seed users
	users, err := seedUsers(db)
	if err != nil {
		return err
	}

	// Seed accounts
	accounts, err := seedAccounts(db, users)
	if err != nil {
		return err
	}

	// Seed transactions
	if err := seedTransactions(db, accounts); err != nil {
		return err
	}

	return nil
}

func clearData(db *gorm.DB) error {
	// Drop tables in reverse order to avoid foreign key constraints
	if err := db.Exec("DELETE FROM transactions").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM accounts").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	return nil
}

func seedUsers(db *gorm.DB) ([]models.User, error) {
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
		if err := db.Create(&users[i]).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

func seedAccounts(db *gorm.DB, users []models.User) ([]models.Account, error) {
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
			}
			
			if err := db.Create(&account).Error; err != nil {
				return nil, err
			}
			
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func seedTransactions(db *gorm.DB, accounts []models.Account) error {
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
		
		if err := db.Create(&initialDeposit).Error; err != nil {
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
			
			if err := db.Create(&transaction).Error; err != nil {
				return err
			}
		}
		
		// Update the account balance to match final transaction
		if err := db.Model(&account).Update("balance", balance).Error; err != nil {
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
				if err := db.Create(&withdrawalTx).Error; err != nil {
					return err
				}
				
				if err := db.Create(&depositTx).Error; err != nil {
					return err
				}
				
				// Update account balances
				if err := db.Model(&fromAccount).Update("balance", fromBalance).Error; err != nil {
					return err
				}
				
				if err := db.Model(&toAccount).Update("balance", toBalance).Error; err != nil {
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
