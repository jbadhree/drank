package functional

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/config"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/handlers"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/middleware"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
	"google.golang.org/api/option"
)

var (
	testFirestoreClient *firestore.Client
	testAuthClient      *auth.Client
	testRouter          *gin.Engine
	testConfig          *config.Config
	testContext         context.Context
)

// SetupTestFirebase initializes Firebase clients for testing
func SetupTestFirebase(t *testing.T) (*firestore.Client, *auth.Client, error) {
	ctx := context.Background()

	// Set emulator environment variables if not already set
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8091")
	}
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") == "" {
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	}

	// Create Firebase app configuration
	config := &firebase.Config{
		ProjectID: "test-project",
	}

	// Create Firebase app
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing Firebase app: %v", err)
	}

	// Get Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing Firebase Auth client: %v", err)
	}

	// Get Firestore client
	firestoreClient, err := firestore.NewClient(ctx, "test-project", option.WithoutAuthentication())
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing Firestore client: %v", err)
	}

	return firestoreClient, authClient, nil
}

// CleanupCollection removes all documents from a collection
func CleanupCollection(client *firestore.Client, collection string) error {
	ctx := context.Background()
	iter := client.Collection(collection).Documents(ctx)
	batch := client.Batch()
	
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		batch.Delete(doc.Ref)
	}
	
	_, err := batch.Commit(ctx)
	return err
}

// CleanupTestFirebase cleans up all collections and closes clients
func CleanupTestFirebase(firestoreClient *firestore.Client) error {
	// Clean up collections
	collections := []string{"users", "accounts", "transactions"}
	for _, collection := range collections {
		if err := CleanupCollection(firestoreClient, collection); err != nil {
			return err
		}
	}
	
	// Close Firestore client
	return firestoreClient.Close()
}

// SetupTestRouter creates a router with all the routes for testing
func SetupTestRouter(firestoreClient *firestore.Client, authClient *auth.Client) *gin.Engine {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Use default test configuration
	cfg := config.New()
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(firestoreClient)
	accountRepo := repository.NewAccountRepository(firestoreClient)
	transactionRepo := repository.NewTransactionRepository(firestoreClient)
	
	// Initialize services
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	
	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	
	// Initialize router
	router := gin.Default()
	
	// API routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes - no auth required
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/register", authHandler.Register)
		
		// User routes - auth required
		users := v1.Group("/users")
		users.Use(authMiddleware.Authenticate())
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/me", userHandler.GetCurrentUser)
		}
		
		// Account routes - auth required
		accounts := v1.Group("/accounts")
		accounts.Use(authMiddleware.Authenticate())
		{
			accounts.GET("", accountHandler.GetAllAccounts)
			accounts.GET("/:id", accountHandler.GetAccountByID)
			accounts.GET("/user/:userId", accountHandler.GetAccountsByUserID)
			accounts.POST("", accountHandler.CreateAccount)
		}
		
		// Transaction routes - auth required
		transactions := v1.Group("/transactions")
		transactions.Use(authMiddleware.Authenticate())
		{
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.GET("/account/:accountId", transactionHandler.GetTransactionsByAccountID)
			transactions.POST("/transfer", transactionHandler.Transfer)
			transactions.POST("/deposit", transactionHandler.CreateDeposit)
			transactions.POST("/withdrawal", transactionHandler.CreateWithdrawal)
		}
	}
	
	return router
}

// SetupTest initializes everything for tests
func SetupTest(t *testing.T) {
	// Initialize Firebase only once
	if testFirestoreClient == nil || testAuthClient == nil {
		var err error
		testFirestoreClient, testAuthClient, err = SetupTestFirebase(t)
		if err != nil {
			t.Fatalf("Failed to set up test Firebase: %v", err)
		}
		
		testContext = context.Background()
	}
	
	// Clean up any existing data
	collections := []string{"users", "accounts", "transactions"}
	for _, collection := range collections {
		if err := CleanupCollection(testFirestoreClient, collection); err != nil {
			t.Logf("Warning: Failed to clean up collection %s: %v", collection, err)
		}
	}
	
	// Initialize router only once
	if testRouter == nil {
		testRouter = SetupTestRouter(testFirestoreClient, testAuthClient)
	}
	
	// Initialize config
	if testConfig == nil {
		testConfig = config.New()
	}
}

// MakeRequest is a helper function to make HTTP requests for tests
func MakeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	
	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	
	return w
}

// CreateTestUser creates a test user for tests
func CreateTestUser(email, password, firstName, lastName string) (models.User, error) {
	// Hash the password
	hashedPassword, err := models.GeneratePasswordHash(password)
	if err != nil {
		return models.User{}, err
	}
	
	// Create user
	user := models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Add to Firestore
	docRef, _, err := testFirestoreClient.Collection("users").Add(testContext, user)
	if err != nil {
		return models.User{}, err
	}
	
	// Update the user with the generated ID
	user.ID = docRef.ID
	_, err = docRef.Set(testContext, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	if err != nil {
		return models.User{}, err
	}
	
	return user, nil
}

// CreateTestAccount creates a test account for tests
func CreateTestAccount(userID string, accountNumber string, accountType models.AccountType, balance float64) (models.Account, error) {
	// Create account
	account := models.Account{
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   accountType,
		Balance:       balance,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	// Add to Firestore
	docRef, _, err := testFirestoreClient.Collection("accounts").Add(testContext, account)
	if err != nil {
		return models.Account{}, err
	}
	
	// Update the account with the generated ID
	account.ID = docRef.ID
	_, err = docRef.Set(testContext, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	if err != nil {
		return models.Account{}, err
	}
	
	return account, nil
}

// LoginTestUser logs in a test user and returns the auth token
func LoginTestUser(email, password string) (string, error) {
	loginReq := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	
	w := MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	
	var response models.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		return "", err
	}
	
	if w.Code != http.StatusOK {
		// Try to get error message
		var errorResp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &errorResp); err == nil {
			return "", fmt.Errorf("login failed: %s", errorResp["error"])
		}
		return "", fmt.Errorf("login failed with status code: %d", w.Code)
	}
	
	return response.Token, nil
}

// TestMain is the main entry point for tests
func TestMain(m *testing.M) {
	// Check if emulators are running
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" && os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") == "" {
		log.Println("Firebase emulators are not running. Please start them before running the tests.")
		log.Println("Setting up environment variables for emulators...")
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8091")
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	}
	
	// Run the tests
	exitCode := m.Run()
	
	// Cleanup after all tests
	if testFirestoreClient != nil {
		if err := CleanupTestFirebase(testFirestoreClient); err != nil {
			fmt.Printf("Failed to clean up test Firebase: %v\n", err)
		}
	}
	
	os.Exit(exitCode)
}
