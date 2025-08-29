package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gbt-be-template/internal/models"
	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req *models.UserCreateRequest) (*models.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) List(ctx context.Context, page, limit int) ([]*models.UserResponse, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) Login(ctx context.Context, req *models.UserLoginRequest) (string, *models.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*models.UserResponse), args.Error(2)
}

func (m *MockUserService) Logout(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func setupUserHandler() (*UserHandler, *MockUserService) {
	mockService := &MockUserService{}
	log := logger.New("info", "text")
	handler := NewUserHandler(mockService, log)
	return handler, mockService
}

func TestUserHandler_Create(t *testing.T) {
	handler, mockService := setupUserHandler()

	t.Run("successful creation", func(t *testing.T) {
		req := &models.UserCreateRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		expectedResponse := &models.UserResponse{
			ID:        1,
			Email:     req.Email,
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		mockService.On("Create", mock.Anything, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Create(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Create(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		req := &models.UserCreateRequest{
			Email: "invalid-email", // Invalid email format
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Create(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("service error", func(t *testing.T) {
		req := &models.UserCreateRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		mockService.On("Create", mock.Anything, req).Return(nil, errors.New("email already exists"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Create(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	handler, mockService := setupUserHandler()

	t.Run("successful get", func(t *testing.T) {
		expectedResponse := &models.UserResponse{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
		}

		mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		recorder := httptest.NewRecorder()

		// Setup chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handler.GetByID(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
		recorder := httptest.NewRecorder()

		// Setup chi context with invalid ID
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handler.GetByID(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))

		request := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		recorder := httptest.NewRecorder()

		// Setup chi context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handler.GetByID(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Login(t *testing.T) {
	handler, mockService := setupUserHandler()

	t.Run("successful login", func(t *testing.T) {
		req := &models.UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := &models.UserResponse{
			ID:    1,
			Email: req.Email,
		}

		mockService.On("Login", mock.Anything, req).Return("token123", expectedUser, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Login(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		
		var response map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "token123", data["access_token"])
		
		mockService.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req := &models.UserLoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockService.On("Login", mock.Anything, req).Return("", nil, errors.New("invalid credentials"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		handler.Login(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Logout(t *testing.T) {
	handler, mockService := setupUserHandler()

	t.Run("successful logout", func(t *testing.T) {
		mockService.On("Logout", mock.Anything, uint(1)).Return(nil)

		request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		recorder := httptest.NewRecorder()

		// Add user ID to context (simulating authenticated user)
		ctx := context.WithValue(request.Context(), middleware.UserIDKey, uint(1))
		request = request.WithContext(ctx)

		handler.Logout(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("user not authenticated", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		recorder := httptest.NewRecorder()

		handler.Logout(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}
