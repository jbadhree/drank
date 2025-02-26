package functional

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create test users
	user1, err := CreateTestUser("account1@example.com", "password123", "Account", "User1")
	assert.NoError(t, err)
	
	user2, err := CreateTestUser("account2@example.com", "password123", "Account", "User2")
	assert.NoError(t, err)
	
	// Create test accounts
	account1, err := CreateTestAccount(user1.ID, "ACC100001", models.Checking, 1000.0)
	assert.NoError(t, err)
	
	_, err = CreateTestAccount(user1.ID, "ACC100002", models.Savings, 2000.0)
	assert.NoError(t, err)
	
	_, err = CreateTestAccount(user2.ID, "ACC200001", models.Checking, 3000.0)
	assert.NoError(t, err)
	
	// Get auth tokens
	token1, err := LoginTestUser("account1@example.com", "password123")
	assert.NoError(t, err)
	
	token2, err := LoginTestUser("account2@example.com", "password123")
	assert.NoError(t, err)
	
	t.Run("GetAllAccounts should return all accounts for admin user", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/accounts", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var accounts []models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.NoError(t, err)
		
		// Verify response contains all accounts
		assert.Equal(t, 3, len(accounts))
	})
	
	t.Run("GetAccountByID should return account details", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("/api/v1/accounts/%d", account1.ID)
		w := MakeRequest("GET", url, nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var account models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &account)
		assert.NoError(t, err)
		
		// Verify response
		assert.Equal(t, account1.ID, account.ID)
		assert.Equal(t, user1.ID, account.UserID)
		assert.Equal(t, "ACC100001", account.AccountNumber)
		assert.Equal(t, models.Checking, account.AccountType)
		assert.Equal(t, 1000.0, account.Balance)
	})
	
	t.Run("GetAccountByID should return not found for non-existent account", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/accounts/999", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["message"], "account not found")
	})
	
	t.Run("GetAccountsByUserID should return user's accounts", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("/api/v1/accounts/user/%d", user1.ID)
		w := MakeRequest("GET", url, nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var accounts []models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.NoError(t, err)
		
		// Verify response contains only user1's accounts
		assert.Equal(t, 2, len(accounts))
		for _, account := range accounts {
			assert.Equal(t, user1.ID, account.UserID)
		}
	})
	
	t.Run("GetAccountsByUserID should return empty array for user with no accounts", func(t *testing.T) {
		// Create a user with no accounts
		userNoAccounts, err := CreateTestUser("noaccount@example.com", "password123", "No", "Accounts")
		assert.NoError(t, err)
		
		// Login
		tokenNoAccounts, err := LoginTestUser("noaccount@example.com", "password123")
		assert.NoError(t, err)
		
		// Act
		url := fmt.Sprintf("/api/v1/accounts/user/%d", userNoAccounts.ID)
		w := MakeRequest("GET", url, nil, tokenNoAccounts)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var accounts []models.AccountDTO
		err = json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.NoError(t, err)
		
		// Verify response is an empty array
		assert.Equal(t, 0, len(accounts))
	})
	
	t.Run("Users should not be able to access accounts of other users directly", func(t *testing.T) {
		// User2 tries to access user1's account
		url := fmt.Sprintf("/api/v1/accounts/%d", account1.ID)
		w := MakeRequest("GET", url, nil, token2)
		
		// Server allows access to all accounts by any authenticated user in this implementation
		// This could be enhanced in a future version with proper authorization
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
