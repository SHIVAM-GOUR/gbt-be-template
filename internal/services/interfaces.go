package services

import (
	"context"

	"gbt-be-template/internal/models"
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(ctx context.Context, req *models.UserCreateRequest) (*models.UserResponse, error)
	GetByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	Update(ctx context.Context, id uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	AdminUpdate(ctx context.Context, id uint, req *models.AdminUserUpdateRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, limit int) ([]*models.UserResponse, int64, error)
	Login(ctx context.Context, req *models.UserLoginRequest) (string, *models.UserResponse, error)
	Logout(ctx context.Context, userID uint) error
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	GenerateToken(userID uint, email string, isAdmin bool) (string, error)
	ValidateToken(token string) (*models.User, error)
	RefreshToken(token string) (string, error)
}

// Services holds all service interfaces
type Services struct {
	User UserService
	Auth AuthService
}
