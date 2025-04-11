package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseClient - Firebase client for auth and firestore
type FirebaseClient struct {
	Auth      *auth.Client
	Firestore *firestore.Client
}

// NewFirebaseClient - Create a new Firebase client
func NewFirebaseClient(cfg *Config) (*FirebaseClient, error) {
	ctx := context.Background()

	// Create Firebase app configuration
	config := &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}

	// Create Firebase app with emulator configuration
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
	firestoreClient, err := firestore.NewClient(ctx, cfg.FirebaseProjectID, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v", err)
		return nil, err
	}

	return &FirebaseClient{
		Auth:      authClient,
		Firestore: firestoreClient,
	}, nil
}

// Close - Close the Firebase client connections
func (fc *FirebaseClient) Close() error {
	return fc.Firestore.Close()
}
