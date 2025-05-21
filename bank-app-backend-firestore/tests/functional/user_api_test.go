package functional

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create test users
	user1, err := CreateTestUser("user1@example.com", "password123", "User", "One")
	assert.NoError(t, err)
	
	user2, err := CreateTestUser("user2@example.com", "password123", "User", "Two")
	assert.NoError(t, err)
	
	// Get auth tokens
	token1, err := LoginTestUser("user1@example.com", "password123")
	assert.NoError(t, err)
	
	token2, err := LoginTestUser("user2@example.com", "password123")
	assert.NoError(t, err)
	
	t.Run("GetAllUsers should return list of users", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var users []models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &users)
		assert.NoError(t, err)
		
		// Verify response contains at least our test users
		assert.GreaterOrEqual(t, len(users), 2)
		
		// Verify user data does not contain passwords
		jsonStr := w.Body.String()
		assert.NotContains(t, jsonStr, "password")
		assert.NotContains(t, jsonStr, "Password")
	})
	
	t.Run("GetUserByID should return user details", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("/api/v1/users/%d", user1.ID)
		w := MakeRequest("GET", url, nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var user models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		
		// Verify response
		assert.Equal(t, user1.ID, user.ID)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "User", user.FirstName)
		assert.Equal(t, "One", user.LastName)
		
		// Verify password is not exposed
		jsonStr := w.Body.String()
		assert.NotContains(t, jsonStr, "password")
		assert.NotContains(t, jsonStr, "Password")
	})
	
	t.Run("GetUserByID should return not found for non-existent user", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users/999", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message contains "User not found"
		assert.Contains(t, response["message"], "User not found")
	})
	
	t.Run("GetCurrentUser should return authenticated user's details", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users/me", nil, token1)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var user models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		
		// Verify response matches the authenticated user
		assert.Equal(t, user1.ID, user.ID)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "User", user.FirstName)
		assert.Equal(t, "One", user.LastName)
	})
	
	t.Run("Different users should get their own details with GetCurrentUser", func(t *testing.T) {
		// Act - user2 requests their own details
		w := MakeRequest("GET", "/api/v1/users/me", nil, token2)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var user models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		
		// Verify response matches user2, not user1
		assert.Equal(t, user2.ID, user.ID)
		assert.Equal(t, "user2@example.com", user.Email)
		assert.Equal(t, "User", user.FirstName)
		assert.Equal(t, "Two", user.LastName)
	})
}
