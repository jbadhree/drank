package functional

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jbadhree/drank/bank-app-backend/internal/config"
	"github.com/jbadhree/drank/bank-app-backend/internal/handlers"
	"github.com/jbadhree/drank/bank-app-backend/internal/middleware"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/repository"
	"github.com/jbadhree/drank/bank-app-backend/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
	testConfig *config.Config
)

// SetupTestDB initializes a test database connection
func SetupTestDB(t *testing.T) (*gorm.DB, error) {
	// Load test environment variables
	cfg := config.New()
	
	// Use test database settings with different port (5435 for test db)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, 5435, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	// Configure the database with minimal logging for tests
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	
	db, err := gorm.Open(postgres.Open(dsn), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %v", err)
	}
	
	// Auto-migrate the schema for test database
	err = db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate test database schema: %v", err)
	}
	
	return db, nil
}

// CleanupTestDB drops all test tables
func CleanupTestDB(db *gorm.DB) error {
	// Get a generic database object
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	// Close the database connection
	return sqlDB.Close()
}

// SetupTestRouter creates a router with all the routes for testing
func SetupTestRouter(db *gorm.DB) *gin.Engine {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Use default test configuration
	cfg := config.New()
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	
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
		}
		
		// Transaction routes - auth required
		transactions := v1.Group("/transactions")
		transactions.Use(authMiddleware.Authenticate())
		{
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.GET("/account/:accountId", transactionHandler.GetTransactionsByAccountID)
			transactions.POST("/transfer", transactionHandler.Transfer)
		}
	}
	
	return router
}

// SetupTest initializes everything for tests
func SetupTest(t *testing.T) {
	// Initialize database only once
	if testDB == nil {
		var err error
		testDB, err = SetupTestDB(t)
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
	}
	
	// Clean up any existing data
	testDB.Exec("TRUNCATE users, accounts, transactions RESTART IDENTITY CASCADE")
	
	// Initialize router only once
	if testRouter == nil {
		testRouter = SetupTestRouter(testDB)
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
func CreateTestUser(email, password, firstName, lastName string) (*models.User, error) {
	user := &models.User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
	
	if err := testDB.Create(user).Error; err != nil {
		return nil, err
	}
	
	return user, nil
}

// CreateTestAccount creates a test account for tests
func CreateTestAccount(userID uint, accountNumber string, accountType models.AccountType, balance float64) (*models.Account, error) {
	account := &models.Account{
		UserID:        userID,
		AccountNumber: accountNumber,
		AccountType:   accountType,
		Balance:       balance,
	}
	
	if err := testDB.Create(account).Error; err != nil {
		return nil, err
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
		return "", fmt.Errorf("login failed with status code: %d", w.Code)
	}
	
	return response.Token, nil
}

// TestMain is the main entry point for tests
func TestMain(m *testing.M) {
	// Run the tests
	exitCode := m.Run()
	
	// Cleanup after all tests
	if testDB != nil {
		if err := CleanupTestDB(testDB); err != nil {
			fmt.Printf("Failed to clean up test database: %v\n", err)
		}
	}
	
	os.Exit(exitCode)
}
