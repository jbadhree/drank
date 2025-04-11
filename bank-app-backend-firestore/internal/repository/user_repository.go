package repository

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserRepositoryImpl - Implementation of the UserRepository interface
type UserRepositoryImpl struct {
	client *firestore.Client
	ctx    context.Context
}

// NewUserRepository - Create a new user repository
func NewUserRepository(client *firestore.Client) interfaces.UserRepository {
	return &UserRepositoryImpl{
		client: client,
		ctx:    context.Background(),
	}
}

// Create - Create a new user
func (r *UserRepositoryImpl) Create(user models.User) (models.User, error) {
	// Check if user already exists
	query := r.client.Collection("users").Where("email", "==", user.Email).Limit(1)
	iter := query.Documents(r.ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err == nil {
		return models.User{}, errors.New("user with this email already exists")
	}
	if err != iterator.Done {
		return models.User{}, err
	}

	// Set created and updated timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Add user to Firestore
	docRef, _, err := r.client.Collection("users").Add(r.ctx, user)
	if err != nil {
		return models.User{}, err
	}

	// Update the user with the generated ID
	user.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{
		"id": docRef.ID,
	}, firestore.MergeAll)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// FindByID - Find user by ID
func (r *UserRepositoryImpl) FindByID(id string) (models.User, error) {
	docRef := r.client.Collection("users").Doc(id)
	docSnapshot, err := docRef.Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	var user models.User
	err = docSnapshot.DataTo(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// FindByEmail - Find user by email
func (r *UserRepositoryImpl) FindByEmail(email string) (models.User, error) {
	query := r.client.Collection("users").Where("email", "==", email).Limit(1)
	iter := query.Documents(r.ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// FindAll - Find all users
func (r *UserRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User

	iter := r.client.Collection("users").Documents(r.ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var user models.User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Update - Update a user
func (r *UserRepositoryImpl) Update(user models.User) (models.User, error) {
	// Check if user exists
	docRef := r.client.Collection("users").Doc(user.ID)
	_, err := docRef.Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	// Update user
	user.UpdatedAt = time.Now()
	_, err = docRef.Set(r.ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Delete - Delete a user
func (r *UserRepositoryImpl) Delete(id string) error {
	_, err := r.client.Collection("users").Doc(id).Delete(r.ctx)
	if err != nil {
		return err
	}

	return nil
}
