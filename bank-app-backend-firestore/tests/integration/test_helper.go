package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// TestFirebaseClient holds Firestore and Auth clients for testing
type TestFirebaseClient struct {
	Auth      *auth.Client
	Firestore *firestore.Client
	ctx       context.Context
}

// SetupTestFirebase initializes Firebase clients pointing to emulators
func SetupTestFirebase() (*TestFirebaseClient, error) {
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
		log.Fatalf("Error initializing Firebase app: %v", err)
		return nil, err
	}

	// Get Firebase Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v", err)
		return nil, err
	}

	// Get Firestore client
	firestoreClient, err := firestore.NewClient(ctx, "test-project", option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
		return nil, err
	}

	return &TestFirebaseClient{
		Auth:      authClient,
		Firestore: firestoreClient,
		ctx:       ctx,
	}, nil
}

// CleanupCollection removes all documents from a collection
func (tfc *TestFirebaseClient) CleanupCollection(collection string) error {
	iter := tfc.Firestore.Collection(collection).Documents(tfc.ctx)
	batch := tfc.Firestore.Batch()
	
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		batch.Delete(doc.Ref)
	}
	
	_, err := batch.Commit(tfc.ctx)
	return err
}

// CleanupAllCollections cleans up all collections used in tests
func (tfc *TestFirebaseClient) CleanupAllCollections() error {
	collections := []string{"users", "accounts", "transactions"}
	
	for _, collection := range collections {
		if err := tfc.CleanupCollection(collection); err != nil {
			return err
		}
	}
	
	return nil
}

// Close closes all connections
func (tfc *TestFirebaseClient) Close() error {
	return tfc.Firestore.Close()
}

// TestMain is the entry point for all tests in this package
func TestMain(m *testing.M) {
	// Setup
	firebase, err := SetupTestFirebase()
	if err != nil {
		log.Fatalf("Failed to set up Firebase for testing: %v", err)
	}
	
	// Clean up any existing data
	if err := firebase.CleanupAllCollections(); err != nil {
		log.Printf("Warning: Failed to clean up collections: %v", err)
	}
	
	// Run tests
	code := m.Run()
	
	// Clean up after tests
	if err := firebase.CleanupAllCollections(); err != nil {
		log.Printf("Warning: Failed to clean up collections: %v", err)
	}
	
	// Close connections
	if err := firebase.Close(); err != nil {
		log.Printf("Warning: Failed to close Firebase client: %v", err)
	}
	
	os.Exit(code)
}
