package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:100"`
	Password  string         `json:"-" gorm:"not null;size:255"` // "-" excludes from JSON
	FirstName string         `json:"first_name" gorm:"size:100"`
	LastName  string         `json:"last_name" gorm:"size:100"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

// UserUpdateRequest represents the request payload for updating a user
type UserUpdateRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint       `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	IsActive  bool       `json:"is_active"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		IsAdmin:   u.IsAdmin,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Any pre-creation logic can go here
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Any pre-update logic can go here
	return nil
}
