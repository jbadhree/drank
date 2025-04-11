package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User - User model for Firestore
type User struct {
	ID        string    `json:"id" firestore:"id"`
	Email     string    `json:"email" firestore:"email"`
	Password  string    `json:"-" firestore:"password"` // Password is not exposed in JSON
	FirstName string    `json:"firstName" firestore:"firstName"`
	LastName  string    `json:"lastName" firestore:"lastName"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}

// GeneratePasswordHash - Generate a hash for the password
func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword - Compare the password with the hashed password
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// ToDTO - Convert User model to DTO (Data Transfer Object)
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserDTO - Data Transfer Object for User
type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginRequest - Request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse - Response body for login
type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
