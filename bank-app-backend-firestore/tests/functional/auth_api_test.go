package functional

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI(t *testing.T) {
	// Set up the test environment
	SetupTest(t)
	
	// Create a test user for authentication tests
	user, err := CreateTestUser("auth@example.com", "password123", "Auth", "User")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	
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
		assert.Contains(t, response["error"], "Invalid credentials")
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
		assert.Contains(t, response["error"], "Invalid credentials")
	})
	
	t.Run("Register a new user should succeed", func(t *testing.T) {
		// Arrange
		registerReq := models.User{
			Email:     "newuser@example.com",
			Password:  "password123",
			FirstName: "New",
			LastName:  "User",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/register", registerReq, "")
		
		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response models.UserDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "newuser@example.com", response.Email)
		assert.Equal(t, "New", response.FirstName)
		assert.Equal(t, "User", response.LastName)
	})
	
	t.Run("Register with existing email should fail", func(t *testing.T) {
		// Arrange
		registerReq := models.User{
			Email:     "auth@example.com", // Already exists
			Password:  "password123",
			FirstName: "Another",
			LastName:  "User",
		}
		
		// Act
		w := MakeRequest("POST", "/api/v1/auth/register", registerReq, "")
		
		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Verify error message about email already existing
		assert.Contains(t, response["error"], "already exists")
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
		assert.NotEmpty(t, response["error"])
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
		assert.NotEmpty(t, response["error"])
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
