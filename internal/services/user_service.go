package services

import (
	"context"
	"errors"
	"fmt"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/models"
	"gbt-be-template/internal/repository"
	"gbt-be-template/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
	authSvc  AuthService
	cfg      *config.Config
	log      *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, authSvc AuthService, cfg *config.Config, log *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		authSvc:  authSvc,
		cfg:      cfg,
		log:      log,
	}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Check if user already exists by email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to check if user exists by email")
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Check if username is taken
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.log.WithError(err).Error("Failed to check if user exists by username")
		return nil, fmt.Errorf("failed to check username availability: %w", err)
	}
	if exists {
		return nil, errors.New("username is already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		IsAdmin:   false,
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.log.WithField("user_id", user.ID).Info("User created successfully")
	return user.ToResponse(), nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

// Update updates a user
func (s *userService) Update(ctx context.Context, id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to get user for update")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Email != nil && *req.Email != user.Email {
		// Check if new email is already taken
		exists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}
		if exists {
			return nil, errors.New("email is already taken")
		}
		user.Email = *req.Email
	}

	if req.Username != nil && *req.Username != user.Username {
		// Check if new username is already taken
		exists, err := s.userRepo.ExistsByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username availability: %w", err)
		}
		if exists {
			return nil, errors.New("username is already taken")
		}
		user.Username = *req.Username
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.log.WithField("user_id", id).Info("User updated successfully")
	return user.ToResponse(), nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uint) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to get user for deletion")
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.log.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// List retrieves a paginated list of users
func (s *userService) List(ctx context.Context, page, limit int) ([]*models.UserResponse, int64, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get users
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		s.log.WithError(err).Error("Failed to list users")
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to count users")
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert to response format
	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(ctx context.Context, req *models.UserLoginRequest) (string, *models.UserResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.log.WithError(err).WithField("email", req.Email).Error("Failed to get user for login")
		return "", nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, errors.New("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.log.WithField("email", req.Email).Warn("Invalid password attempt")
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.authSvc.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		s.log.WithError(err).WithField("user_id", user.ID).Error("Failed to generate token")
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.log.WithError(err).WithField("user_id", user.ID).Warn("Failed to update last login")
	}

	s.log.WithField("user_id", user.ID).Info("User logged in successfully")
	return token, user.ToResponse(), nil
}

// Logout logs out a user (in this implementation, it's just a placeholder)
func (s *userService) Logout(ctx context.Context, userID uint) error {
	// In a more sophisticated implementation, you might:
	// - Add the token to a blacklist
	// - Store logout time in database
	// - Invalidate refresh tokens
	
	s.log.WithField("user_id", userID).Info("User logged out successfully")
	return nil
}
