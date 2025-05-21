package repository

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type userRepository struct {
	client *firestore.Client
	ctx    context.Context
}

func NewUserRepository(client *firestore.Client) UserRepository {
	return &userRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *userRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	docRef, _, err := r.client.Collection("users").Add(r.ctx, user)
	if err != nil {
		return err
	}
	user.ID = docRef.ID
	_, err = docRef.Set(r.ctx, map[string]interface{}{"id": docRef.ID}, firestore.MergeAll)
	return err
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	doc, err := r.client.Collection("users").Doc(id).Get(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	iter := r.client.Collection("users").Where("email", "==", email).Limit(1).Documents(r.ctx)
	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	iter := r.client.Collection("users").Documents(r.ctx)
	var users []models.User
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

func (r *userRepository) Update(user *models.User) error {
	if user.ID == "" {
		return errors.New("user ID is required for update")
	}
	user.UpdatedAt = time.Now()
	_, err := r.client.Collection("users").Doc(user.ID).Set(r.ctx, user)
	return err
}

func (r *userRepository) Delete(id string) error {
	_, err := r.client.Collection("users").Doc(id).Delete(r.ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errors.New("user not found")
		}
		return err
	}
	return err
}
