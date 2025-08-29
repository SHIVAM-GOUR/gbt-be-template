package repository

import (
	"context"
	"testing"
	"time"

	"gbt-be-template/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *Database {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	database := &Database{DB: db}

	// Auto migrate
	err = database.AutoMigrate()
	require.NoError(t, err)

	return database
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test getting the user
	foundUser, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Username, foundUser.Username)

	// Test getting non-existent user
	notFoundUser, err := repo.GetByID(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, notFoundUser)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test getting the user by email
	foundUser, err := repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test getting non-existent user
	notFoundUser, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFoundUser)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update the user
	user.FirstName = "Updated"
	user.LastName = "Name"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify the update
	updatedUser, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updatedUser.FirstName)
	assert.Equal(t, "Name", updatedUser.LastName)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	// Verify the user is soft deleted
	deletedUser, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedUser) // Should be nil due to soft delete
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Test existing email
	exists, err := repo.ExistsByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test non-existing email
	exists, err = repo.ExistsByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &models.User{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
		IsAdmin:   false,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update last login
	err = repo.UpdateLastLogin(ctx, user.ID)
	assert.NoError(t, err)

	// Verify last login was updated
	updatedUser, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser.LastLogin)
	assert.WithinDuration(t, time.Now(), *updatedUser.LastLogin, time.Minute)
}
