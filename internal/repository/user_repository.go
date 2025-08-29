package repository

import (
	"context"
	"errors"
	"time"

	"gbt-be-template/internal/models"

	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface
type userRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.DB.WithContext(ctx).Save(user).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.DB.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List retrieves a list of users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := r.db.DB.WithContext(ctx).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	
	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.DB.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUsername checks if a user exists with the given username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.DB.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateLastLogin updates the last login time for a user
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("last_login", now).Error
}
