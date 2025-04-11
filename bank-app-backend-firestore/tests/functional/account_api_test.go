package functional

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccountAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create test users
	user1, err := CreateTestUser("account_test_1@example.com", "password123", "Account", "TestOne")
	assert.NoError(t, err)
	
	user2, err := CreateTestUser("account_test_2@example.com", "password123", "Account", "TestTwo")
	assert.NoError(t, err)
	
	// Create test accounts
	account1, err := CreateTestAccount(user1.ID, "1000000001", models.Checking, 1000.00)
	assert.NoError(t, err)
	
	account2, err := CreateTestAccount(user1.ID, "1000000002", models.Savings, 2000.00)
	assert.NoError(t, err)
	
	account3, err := CreateTestAccount(user2.ID, "1000000003", models.Checking, 3000.00)
	assert.NoError(t, err)
	
	// Get tokens for authentication
	token1, err := LoginTestUser("account_test_1@example.com", "password123")
	assert.NoError(t, err)
	
	token2, err := LoginTestUser("account_test_2@example.com", "password123")
	assert.NoError(t, err)
	
	t.Run("Get all accounts should return a list of accounts", func(t *testing.T) {
		// Act - As admin or user with permissions
		w := MakeRequest("GET", "/api/v1/accounts", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var accounts []models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.NoError(t, err)
		
		// Verify at least 3 accounts are returned
		assert.GreaterOrEqual(t, len(accounts), 3)
		
		// Verify our test accounts are in the list
		foundAccount1 := false
		foundAccount2 := false
		foundAccount3 := false
		for _, a := range accounts {
			if a.ID == account1.ID {
				foundAccount1 = true
				assert.Equal(t, user1.ID, a.UserID)
				assert.Equal(t, "1000000001", a.AccountNumber)
			}
			if a.ID == account2.ID {
				foundAccount2 = true
				assert.Equal(t, user1.ID, a.UserID)
				assert.Equal(t, "1000000002", a.AccountNumber)
			}
			if a.ID == account3.ID {
				foundAccount3 = true
				assert.Equal(t, user2.ID, a.UserID)
				assert.Equal(t, "1000000003", a.AccountNumber)
			}
		}
		assert.True(t, foundAccount1, "Account1 should be in the list")
		assert.True(t, foundAccount2, "Account2 should be in the list")
		assert.True(t, foundAccount3, "Account3 should be in the list")
	})
	
	t.Run("Get account by ID should return the correct account", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/accounts/"+account1.ID, nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var returnedAccount models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &returnedAccount)
		assert.NoError(t, err)
		
		// Verify the correct account is returned
		assert.Equal(t, account1.ID, returnedAccount.ID)
		assert.Equal(t, user1.ID, returnedAccount.UserID)
		assert.Equal(t, "1000000001", returnedAccount.AccountNumber)
		assert.Equal(t, models.Checking, returnedAccount.AccountType)
		assert.Equal(t, 1000.00, returnedAccount.Balance)
	})
	
	t.Run("Get account by ID with invalid ID should return 404", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/accounts/nonexistent-id", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "not found")
	})
	
	t.Run("Get accounts by user ID should return only that user's accounts", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/accounts/user/"+user1.ID, nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var accounts []models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &accounts)
		assert.NoError(t, err)
		
		// Verify only user1's accounts are returned
		assert.Equal(t, 2, len(accounts))
		
		// Verify the correct accounts are returned
		foundAccount1 := false
		foundAccount2 := false
		for _, a := range accounts {
			assert.Equal(t, user1.ID, a.UserID, "All accounts should belong to user1")
			
			if a.ID == account1.ID {
				foundAccount1 = true
				assert.Equal(t, "1000000001", a.AccountNumber)
			}
			if a.ID == account2.ID {
				foundAccount2 = true
				assert.Equal(t, "1000000002", a.AccountNumber)
			}
		}
		assert.True(t, foundAccount1, "Account1 should be in the list")
		assert.True(t, foundAccount2, "Account2 should be in the list")
	})
	
	t.Run("User cannot access another user's accounts by user ID", func(t *testing.T) {
		// Act - user2 trying to access user1's accounts
		w := MakeRequest("GET", "/api/v1/accounts/user/"+user1.ID, nil, token2)
		
		// Assert - should be forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "Access denied")
	})
	
	t.Run("Create new account should succeed", func(t *testing.T) {
		// Arrange
		newAccount := models.Account{
			AccountType: models.Savings,
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/accounts", newAccount, token1)
		
		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var createdAccount models.AccountDTO
		err := json.Unmarshal(w.Body.Bytes(), &createdAccount)
		assert.NoError(t, err)
		
		// Verify the account is created correctly
		assert.NotEmpty(t, createdAccount.ID)
		assert.Equal(t, user1.ID, createdAccount.UserID)
		assert.NotEmpty(t, createdAccount.AccountNumber)
		assert.Equal(t, models.Savings, createdAccount.AccountType)
		assert.Equal(t, 0.0, createdAccount.Balance)
		
		// Verify we can retrieve the account
		w = MakeRequest("GET", "/api/v1/accounts/"+createdAccount.ID, nil, token1)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
