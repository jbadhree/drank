package integration

import (
	"testing"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Integration(t *testing.T) {
	// Skip tests if we're not in an integration test environment
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up Firebase client for testing
	firebaseClient, err := SetupTestFirebase()
	require.NoError(t, err)
	defer firebaseClient.Close()

	// Create user repository
	userRepo := repository.NewUserRepository(firebaseClient.Firestore)

	// Clean up test data before starting
	err = firebaseClient.CleanupCollection("users")
	require.NoError(t, err)

	t.Run("Create and FindByID", func(t *testing.T) {
		// Create a test user
		hashedPassword, err := models.GeneratePasswordHash("password123")
		require.NoError(t, err)

		user := models.User{
			Email:     "test@example.com",
			Password:  hashedPassword,
			FirstName: "Test",
			LastName:  "User",
		}

		// Create the user in Firestore
		createdUser, err := userRepo.Create(user)
		require.NoError(t, err)
		assert.NotEmpty(t, createdUser.ID)
		assert.Equal(t, "test@example.com", createdUser.Email)
		assert.NotZero(t, createdUser.CreatedAt)
		assert.NotZero(t, createdUser.UpdatedAt)

		// Find the user by ID
		foundUser, err := userRepo.FindByID(createdUser.ID)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, foundUser.ID)
		assert.Equal(t, "test@example.com", foundUser.Email)
	})

	t.Run("FindByEmail", func(t *testing.T) {
		// Find the user by email
		foundUser, err := userRepo.FindByEmail("test@example.com")
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", foundUser.Email)
	})

	t.Run("Update", func(t *testing.T) {
		// Find the user to update
		user, err := userRepo.FindByEmail("test@example.com")
		require.NoError(t, err)

		// Update user data
		user.FirstName = "Updated"
		user.LastName = "Name"
		updatedUser, err := userRepo.Update(user)
		require.NoError(t, err)

		// Verify the update
		assert.Equal(t, "Updated", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName)
		assert.True(t, updatedUser.UpdatedAt.After(user.UpdatedAt) || updatedUser.UpdatedAt.Equal(user.UpdatedAt))

		// Retrieve again to confirm persistence
		retrievedUser, err := userRepo.FindByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", retrievedUser.FirstName)
		assert.Equal(t, "Name", retrievedUser.LastName)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Create a second user
		hashedPassword, err := models.GeneratePasswordHash("password456")
		require.NoError(t, err)

		user2 := models.User{
			Email:     "another@example.com",
			Password:  hashedPassword,
			FirstName: "Another",
			LastName:  "User",
		}

		// Create the second user
		_, err = userRepo.Create(user2)
		require.NoError(t, err)

		// Retrieve all users
		users, err := userRepo.FindAll()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 2)

		// Verify we can find both test users
		foundFirst := false
		foundSecond := false
		for _, user := range users {
			if user.Email == "test@example.com" {
				foundFirst = true
			}
			if user.Email == "another@example.com" {
				foundSecond = true
			}
		}
		assert.True(t, foundFirst, "First test user should be in results")
		assert.True(t, foundSecond, "Second test user should be in results")
	})

	t.Run("Delete", func(t *testing.T) {
		// Find a user to delete
		user, err := userRepo.FindByEmail("another@example.com")
		require.NoError(t, err)

		// Delete the user
		err = userRepo.Delete(user.ID)
		require.NoError(t, err)

		// Try to find the deleted user
		_, err = userRepo.FindByID(user.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	// Clean up test data after tests
	err = firebaseClient.CleanupCollection("users")
	require.NoError(t, err)
}
