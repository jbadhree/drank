package functional

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create a test user for authentication tests
	_, err := CreateTestUser("auth@example.com", "password123", "Auth", "User")
	assert.NoError(t, err)
	
	t.Run("Login with valid credentials should succeed", func(t *testing.T) {
		// Arrange
		loginReq := models.LoginRequest{
			Email:    "auth@example.com",
			Password: "password123",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify response
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "auth@example.com", response.User.Email)
		assert.Equal(t, "Auth", response.User.FirstName)
		assert.Equal(t, "User", response.User.LastName)
	})
	
	t.Run("Login with invalid email should fail", func(t *testing.T) {
		// Arrange
		loginReq := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["message"], "invalid email or password")
	})
	
	t.Run("Login with incorrect password should fail", func(t *testing.T) {
		// Arrange
		loginReq := models.LoginRequest{
			Email:    "auth@example.com",
			Password: "wrongpassword",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message
		assert.Contains(t, response["message"], "invalid email or password")
	})
	
	t.Run("Login with invalid format should fail", func(t *testing.T) {
		// Arrange - Missing password
		loginReq := map[string]string{
			"email": "auth@example.com",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
		
		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify there is an error message
		assert.Contains(t, response, "message")
	})
	
	t.Run("Access protected routes without token should fail", func(t *testing.T) {
		// Act - try to access a protected route
		w := MakeRequest("GET", "/api/v1/users/me", nil, "")
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify some error message is returned
		assert.NotEmpty(t, response["message"])
	})
	
	t.Run("Access protected routes with invalid token should fail", func(t *testing.T) {
		// Act - try to access a protected route with an invalid token
		w := MakeRequest("GET", "/api/v1/users/me", nil, "invalid-token")
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify some error message is returned
		assert.NotEmpty(t, response["message"])
	})
	
	t.Run("Access protected routes with valid token should succeed", func(t *testing.T) {
		// Arrange - login to get a valid token
		token, err := LoginTestUser("auth@example.com", "password123")
		assert.NoError(t, err)
		
		// Act - access a protected route with the token
		w := MakeRequest("GET", "/api/v1/users/me", nil, token)
		
		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify response
		assert.Equal(t, "auth@example.com", response.Email)
		assert.Equal(t, "Auth", response.FirstName)
		assert.Equal(t, "User", response.LastName)
	})
}
