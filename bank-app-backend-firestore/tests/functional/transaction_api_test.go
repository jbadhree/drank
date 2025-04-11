package functional

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactionAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create test user
	user, err := CreateTestUser("transaction_test@example.com", "password123", "Transaction", "Test")
	assert.NoError(t, err)
	
	// Create test accounts
	checkingAccount, err := CreateTestAccount(user.ID, "2000000001", models.Checking, 1000.00)
	assert.NoError(t, err)
	
	savingsAccount, err := CreateTestAccount(user.ID, "2000000002", models.Savings, 2000.00)
	assert.NoError(t, err)
	
	// Get token for authentication
	token, err := LoginTestUser("transaction_test@example.com", "password123")
	assert.NoError(t, err)
	
	// Create initial transaction directly in Firestore
	transaction := models.Transaction{
		AccountID:       checkingAccount.ID,
		Amount:          500.00,
		Balance:         1500.00, // New balance after this transaction
		Type:            models.Deposit,
		Description:     "Initial deposit",
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	docRef, _, err := testFirestoreClient.Collection("transactions").Add(testContext, transaction)
	assert.NoError(t, err)
	
	transaction.ID = docRef.ID
	_, err = docRef.Set(testContext, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	assert.NoError(t, err)
	
	t.Run("Get all transactions should return a list of transactions", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/transactions", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var transactions []models.TransactionDTO
		err := json.Unmarshal(w.Body.Bytes(), &transactions)
		assert.NoError(t, err)
		
		// Verify at least 1 transaction is returned
		assert.GreaterOrEqual(t, len(transactions), 1)
		
		// Verify our test transaction is in the list
		found := false
		for _, tx := range transactions {
			if tx.ID == transaction.ID {
				found = true
				assert.Equal(t, checkingAccount.ID, tx.AccountID)
				assert.Equal(t, 500.00, tx.Amount)
				assert.Equal(t, 1500.00, tx.Balance)
				assert.Equal(t, models.Deposit, tx.Type)
				assert.Equal(t, "Initial deposit", tx.Description)
			}
		}
		assert.True(t, found, "Test transaction should be in the list")
	})
	
	t.Run("Get transaction by ID should return the correct transaction", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/transactions/"+transaction.ID, nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var returnedTransaction models.TransactionDTO
		err := json.Unmarshal(w.Body.Bytes(), &returnedTransaction)
		assert.NoError(t, err)
		
		// Verify the correct transaction is returned
		assert.Equal(t, transaction.ID, returnedTransaction.ID)
		assert.Equal(t, checkingAccount.ID, returnedTransaction.AccountID)
		assert.Equal(t, 500.00, returnedTransaction.Amount)
		assert.Equal(t, 1500.00, returnedTransaction.Balance)
		assert.Equal(t, models.Deposit, returnedTransaction.Type)
		assert.Equal(t, "Initial deposit", returnedTransaction.Description)
	})
	
	t.Run("Get transaction by ID with invalid ID should return 404", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/transactions/nonexistent-id", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "not found")
	})
	
	t.Run("Get transactions by account ID should return only that account's transactions", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/transactions/account/"+checkingAccount.ID, nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var transactions []models.TransactionDTO
		err := json.Unmarshal(w.Body.Bytes(), &transactions)
		assert.NoError(t, err)
		
		// Verify transactions are returned
		assert.GreaterOrEqual(t, len(transactions), 1)
		
		// Verify all transactions belong to the checking account
		for _, tx := range transactions {
			assert.Equal(t, checkingAccount.ID, tx.AccountID, "All transactions should belong to the checking account")
		}
	})
	
	t.Run("Create deposit transaction should succeed", func(t *testing.T) {
		// Arrange
		depositReq := models.Transaction{
			AccountID:   checkingAccount.ID,
			Amount:      200.00,
			Type:        models.Deposit,
			Description: "Test deposit",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/transactions/deposit", depositReq, token)
		
		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var createdTransaction models.TransactionDTO
		err := json.Unmarshal(w.Body.Bytes(), &createdTransaction)
		assert.NoError(t, err)
		
		// Verify the transaction is created correctly
		assert.NotEmpty(t, createdTransaction.ID)
		assert.Equal(t, checkingAccount.ID, createdTransaction.AccountID)
		assert.Equal(t, 200.00, createdTransaction.Amount)
		assert.Equal(t, models.Deposit, createdTransaction.Type)
		assert.Equal(t, "Test deposit", createdTransaction.Description)
		
		// Balance should have increased by the deposit amount
		// Since initial balance was 1500 after the first deposit
		assert.Equal(t, 1700.00, createdTransaction.Balance)
	})
	
	t.Run("Create withdrawal transaction should succeed", func(t *testing.T) {
		// Arrange
		withdrawalReq := models.Transaction{
			AccountID:   checkingAccount.ID,
			Amount:      100.00, // This will be made negative by the handler
			Type:        models.Withdrawal,
			Description: "Test withdrawal",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/transactions/withdrawal", withdrawalReq, token)
		
		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var createdTransaction models.TransactionDTO
		err := json.Unmarshal(w.Body.Bytes(), &createdTransaction)
		assert.NoError(t, err)
		
		// Verify the transaction is created correctly
		assert.NotEmpty(t, createdTransaction.ID)
		assert.Equal(t, checkingAccount.ID, createdTransaction.AccountID)
		assert.Equal(t, -100.00, createdTransaction.Amount) // Amount should be negative for withdrawals
		assert.Equal(t, models.Withdrawal, createdTransaction.Type)
		assert.Equal(t, "Test withdrawal", createdTransaction.Description)
		
		// Balance should have decreased by the withdrawal amount
		// Since balance was 1700 after the deposit
		assert.Equal(t, 1600.00, createdTransaction.Balance)
	})
	
	t.Run("Transfer between accounts should succeed", func(t *testing.T) {
		// Arrange
		transferReq := models.TransferRequest{
			FromAccountID: checkingAccount.ID,
			ToAccountID:   savingsAccount.ID,
			Amount:        300.00,
			Description:   "Test transfer",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/transactions/transfer", transferReq, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify success message
		assert.Contains(t, response["message"], "successful")
		
		// Verify the source account balance decreased
		w = MakeRequest("GET", "/api/v1/accounts/"+checkingAccount.ID, nil, token)
		assert.Equal(t, http.StatusOK, w.Code)
		
		var sourceAccount models.AccountDTO
		err = json.Unmarshal(w.Body.Bytes(), &sourceAccount)
		assert.NoError(t, err)
		
		assert.Equal(t, 1300.00, sourceAccount.Balance) // 1600 - 300
		
		// Verify the target account balance increased
		w = MakeRequest("GET", "/api/v1/accounts/"+savingsAccount.ID, nil, token)
		assert.Equal(t, http.StatusOK, w.Code)
		
		var targetAccount models.AccountDTO
		err = json.Unmarshal(w.Body.Bytes(), &targetAccount)
		assert.NoError(t, err)
		
		assert.Equal(t, 2300.00, targetAccount.Balance) // 2000 + 300
	})
	
	t.Run("Transfer with insufficient funds should fail", func(t *testing.T) {
		// Arrange
		transferReq := models.TransferRequest{
			FromAccountID: checkingAccount.ID,
			ToAccountID:   savingsAccount.ID,
			Amount:        2000.00, // More than available balance
			Description:   "Test transfer with insufficient funds",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/transactions/transfer", transferReq, token)
		
		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "insufficient")
	})
	
	t.Run("Transfer to same account should fail", func(t *testing.T) {
		// Arrange
		transferReq := models.TransferRequest{
			FromAccountID: checkingAccount.ID,
			ToAccountID:   checkingAccount.ID, // Same account
			Amount:        100.00,
			Description:   "Test transfer to same account",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/transactions/transfer", transferReq, token)
		
		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "same account")
	})
}
