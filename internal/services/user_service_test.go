package services

import (
	"context"
	"errors"
	"testing"

	"gbt-be-template/internal/config"
	"gbt-be-template/internal/models"
	"gbt-be-template/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(userID uint, email string, isAdmin bool) (string, error) {
	args := m.Called(userID, email, isAdmin)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func setupUserService() (*userService, *MockUserRepository, *MockAuthService) {
	mockRepo := &MockUserRepository{}
	mockAuth := &MockAuthService{}
	cfg := &config.Config{}
	log := logger.New("info", "text")
	
	service := &userService{
		userRepo: mockRepo,
		authSvc:  mockAuth,
		cfg:      cfg,
		log:      log,
	}
	
	return service, mockRepo, mockAuth
}

func TestUserService_Create(t *testing.T) {
	service, mockRepo, _ := setupUserService()
	ctx := context.Background()

	req := &models.UserCreateRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	t.Run("successful creation", func(t *testing.T) {
		mockRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
		mockRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(1).(*models.User)
			user.ID = 1 // Simulate database setting ID
		})

		result, err := service.Create(ctx, req)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.Email)
		assert.Equal(t, req.Username, result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo.On("ExistsByEmail", ctx, req.Email).Return(true, nil)

		result, err := service.Create(ctx, req)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	service, mockRepo, mockAuth := setupUserService()
	ctx := context.Background()

	req := &models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &models.User{
		ID:       1,
		Email:    req.Email,
		Password: string(hashedPassword),
		IsActive: true,
		IsAdmin:  false,
	}

	t.Run("successful login", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
		mockAuth.On("GenerateToken", user.ID, user.Email, user.IsAdmin).Return("token123", nil)
		mockRepo.On("UpdateLastLogin", ctx, user.ID).Return(nil)

		token, userResp, err := service.Login(ctx, req)
		
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
		assert.NotNil(t, userResp)
		assert.Equal(t, user.Email, userResp.Email)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)

		token, userResp, err := service.Login(ctx, req)
		
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, userResp)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("inactive user", func(t *testing.T) {
		inactiveUser := *user
		inactiveUser.IsActive = false
		mockRepo.On("GetByEmail", ctx, req.Email).Return(&inactiveUser, nil)

		token, userResp, err := service.Login(ctx, req)
		
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, userResp)
		assert.Contains(t, err.Error(), "deactivated")
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		wrongReq := &models.UserLoginRequest{
			Email:    req.Email,
			Password: "wrongpassword",
		}
		mockRepo.On("GetByEmail", ctx, wrongReq.Email).Return(user, nil)

		token, userResp, err := service.Login(ctx, wrongReq)
		
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, userResp)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	service, mockRepo, _ := setupUserService()
	ctx := context.Background()

	user := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	t.Run("user found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(1)).Return(user, nil)

		result, err := service.GetByID(ctx, 1)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

		result, err := service.GetByID(ctx, 999)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("database error"))

		result, err := service.GetByID(ctx, 1)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
