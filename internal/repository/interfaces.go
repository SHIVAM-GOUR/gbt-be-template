package repository

import (
	"context"

	"gbt-be-template/internal/models"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// Repositories holds all repository interfaces
type Repositories struct {
	User UserRepository
}

// NewRepositories creates a new instance of all repositories
func NewRepositories(db *Database) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}
