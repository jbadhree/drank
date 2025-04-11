package services

import (
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/models"
	"github.com/jbadhree/drank/bank-app-backend-firestore/internal/repository/interfaces"
)

// UserService - Service for user operations
type UserService struct {
	repo interfaces.UserRepository
}

// NewUserService - Create a new user service
func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Create - Create a new user
func (s *UserService) Create(user models.User) (models.UserDTO, error) {
	// Hash the password
	hashedPassword, err := models.GeneratePasswordHash(user.Password)
	if err != nil {
		return models.UserDTO{}, err
	}
	user.Password = hashedPassword

	// Create the user
	createdUser, err := s.repo.Create(user)
	if err != nil {
		return models.UserDTO{}, err
	}

	return createdUser.ToDTO(), nil
}

// GetByID - Get user by ID
func (s *UserService) GetByID(id string) (models.UserDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.ToDTO(), nil
}

// GetByEmail - Get user by email
func (s *UserService) GetByEmail(email string) (models.User, error) {
	return s.repo.FindByEmail(email)
}

// GetAll - Get all users
func (s *UserService) GetAll() ([]models.UserDTO, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	userDTOs := make([]models.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToDTO()
	}

	return userDTOs, nil
}

// Update - Update a user
func (s *UserService) Update(user models.User) (models.UserDTO, error) {
	// Check if password needs to be updated
	if user.Password != "" {
		hashedPassword, err := models.GeneratePasswordHash(user.Password)
		if err != nil {
			return models.UserDTO{}, err
		}
		user.Password = hashedPassword
	}

	// Update the user
	updatedUser, err := s.repo.Update(user)
	if err != nil {
		return models.UserDTO{}, err
	}

	return updatedUser.ToDTO(), nil
}

// Delete - Delete a user
func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}

// Authenticate - Authenticate a user
func (s *UserService) Authenticate(email, password string) (models.User, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	// Check password
	if err := user.ComparePassword(password); err != nil {
		return models.User{}, err
	}

	return user, nil
}
