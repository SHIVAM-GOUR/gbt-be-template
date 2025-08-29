package services

import (
	"fmt"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/models"
	"gbt-be-template/internal/repository"
	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/utils"
)

// authService implements the AuthService interface
type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	log      *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config, log *logger.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
		log:      log,
	}
}

// GenerateToken generates a JWT token for a user
func (s *authService) GenerateToken(userID uint, email string, isAdmin bool) (string, error) {
	token, err := utils.GenerateJWT(userID, email, isAdmin, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		s.log.WithError(err).WithField("user_id", userID).Error("Failed to generate JWT token")
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.log.WithField("user_id", userID).Info("JWT token generated successfully")
	return token, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(token string) (*models.User, error) {
	claims, err := utils.ValidateJWT(token, s.cfg.JWT.Secret)
	if err != nil {
		s.log.WithError(err).Warn("Failed to validate JWT token")
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database to ensure they still exist and are active
	user, err := s.userRepo.GetByID(nil, claims.UserID)
	if err != nil {
		s.log.WithError(err).WithField("user_id", claims.UserID).Error("Failed to get user for token validation")
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}

	if user == nil {
		s.log.WithField("user_id", claims.UserID).Warn("Token validation failed: user not found")
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		s.log.WithField("user_id", claims.UserID).Warn("Token validation failed: user is inactive")
		return nil, fmt.Errorf("user account is deactivated")
	}

	return user, nil
}

// RefreshToken generates a new token with extended expiry
func (s *authService) RefreshToken(token string) (string, error) {
	newToken, err := utils.RefreshJWT(token, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		s.log.WithError(err).Warn("Failed to refresh JWT token")
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	s.log.Info("JWT token refreshed successfully")
	return newToken, nil
}
