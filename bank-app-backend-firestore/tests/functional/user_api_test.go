package functional

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
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
	
	// Get a valid token for authentication
	token, err := LoginTestUser("user1@example.com", "password123")
	assert.NoError(t, err)
	
	t.Run("Get all users should return a list of users", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var users []models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &users)
		assert.NoError(t, err)
		
		// Verify at least 2 users are returned
		assert.GreaterOrEqual(t, len(users), 2)
		
		// Verify our test users are in the list
		foundUser1 := false
		foundUser2 := false
		for _, u := range users {
			if u.Email == "user1@example.com" {
				foundUser1 = true
				assert.Equal(t, "User", u.FirstName)
				assert.Equal(t, "One", u.LastName)
			}
			if u.Email == "user2@example.com" {
				foundUser2 = true
				assert.Equal(t, "User", u.FirstName)
				assert.Equal(t, "Two", u.LastName)
			}
		}
		assert.True(t, foundUser1, "User1 should be in the list")
		assert.True(t, foundUser2, "User2 should be in the list")
	})
	
	t.Run("Get user by ID should return the correct user", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users/"+user2.ID, nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var returnedUser models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &returnedUser)
		assert.NoError(t, err)
		
		// Verify the correct user is returned
		assert.Equal(t, user2.ID, returnedUser.ID)
		assert.Equal(t, "user2@example.com", returnedUser.Email)
		assert.Equal(t, "User", returnedUser.FirstName)
		assert.Equal(t, "Two", returnedUser.LastName)
	})
	
	t.Run("Get user by ID with invalid ID should return 404", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users/nonexistent-id", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["error"], "not found")
	})
	
	t.Run("Get current user should return the authenticated user", func(t *testing.T) {
		// Act
		w := MakeRequest("GET", "/api/v1/users/me", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var returnedUser models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &returnedUser)
		assert.NoError(t, err)
		
		// Verify the authenticated user is returned
		assert.Equal(t, user1.ID, returnedUser.ID)
		assert.Equal(t, "user1@example.com", returnedUser.Email)
		assert.Equal(t, "User", returnedUser.FirstName)
		assert.Equal(t, "One", returnedUser.LastName)
	})
	
	t.Run("Get current user with different user's token", func(t *testing.T) {
		// Arrange - get token for user2
		token2, err := LoginTestUser("user2@example.com", "password123")
		assert.NoError(t, err)
		
		// Act
		w := MakeRequest("GET", "/api/v1/users/me", nil, token2)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var returnedUser models.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &returnedUser)
		assert.NoError(t, err)
		
		// Verify the authenticated user is user2
		assert.Equal(t, user2.ID, returnedUser.ID)
		assert.Equal(t, "user2@example.com", returnedUser.Email)
		assert.Equal(t, "User", returnedUser.FirstName)
		assert.Equal(t, "Two", returnedUser.LastName)
	})
}
