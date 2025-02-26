package services

import (
	"errors"

	"github.com/jbadhree/drank/bank-app-backend/internal/models"
	"github.com/jbadhree/drank/bank-app-backend/internal/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	AuthenticateUser(email, password string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) CreateUser(user *models.User) error {
	// Check if user with the same email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) UpdateUser(user *models.User) error {
	// Check if the user exists
	_, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return err
	}

	// Check if the new email is already taken by another user
	if user.Email != "" {
		existingUser, err := s.userRepo.FindByEmail(user.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return errors.New("email is already taken by another user")
		}
	}

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	// Check if the user exists
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}

func (s *userService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := user.ComparePassword(password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
