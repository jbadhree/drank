package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/config"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/handlers"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/middleware"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/services"
	"github.com/jbadhree/drank/bank-app-backend-firestore/seed"
)

// @title           Banking API (Firestore)
// @version         1.0
// @description     A demo banking application API using Firestore
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token for authentication
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Configure the application
	cfg := config.New()

	// Initialize Firebase client
	firebase, err := config.NewFirebaseClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer firebase.Close()

	// Check if seed flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		log.Println("Seeding database...")
		if err := seed.SeedDatabase(firebase.Firestore); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeded successfully")
		return
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(firebase.Firestore)
	accountRepo := repository.NewAccountRepository(firebase.Firestore)
	transactionRepo := repository.NewTransactionRepository(firebase.Firestore)

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

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS - allow requests from both localhost and the actual server hostname
	// Get frontend URL from environment or use default
	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{"http://localhost:3000"}
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
