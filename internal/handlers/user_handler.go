package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gbt-be-template/internal/models"
	"gbt-be-template/internal/services"
	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/middleware"
	"gbt-be-template/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
	log         *logger.Logger
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
		validator:   validator.New(),
	}
}

// Create handles POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("Invalid JSON in create user request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.log.WithError(err).Warn("Validation failed for create user request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Create user
	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to create user")
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "User created successfully", user)
}

// GetByID handles GET /users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(r.Context(), uint(id))
	if err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("Failed to get user")
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// Update handles PUT /users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	// Check if user is updating their own profile or is admin
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	isAdmin, _ := middleware.GetIsAdminFromContext(r.Context())
	
	if userID != uint(id) && !isAdmin {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You can only update your own profile", nil)
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("Invalid JSON in update user request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.log.WithError(err).Warn("Validation failed for update user request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Update user
	user, err := h.userService.Update(r.Context(), uint(id), &req)
	if err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("Failed to update user")
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "User updated successfully", user)
}

// Delete handles DELETE /users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	// Check if user is deleting their own profile or is admin
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	isAdmin, _ := middleware.GetIsAdminFromContext(r.Context())
	
	if userID != uint(id) && !isAdmin {
		utils.WriteErrorResponse(w, http.StatusForbidden, "You can only delete your own profile", nil)
		return
	}

	if err := h.userService.Delete(r.Context(), uint(id)); err != nil {
		h.log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// List handles GET /users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	users, total, err := h.userService.List(r.Context(), page, limit)
	if err != nil {
		h.log.WithError(err).Error("Failed to list users")
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users", nil)
		return
	}

	utils.WritePaginatedResponse(w, http.StatusOK, "Users retrieved successfully", users, total, page, limit)
}

// Login handles POST /auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Warn("Invalid JSON in login request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.log.WithError(err).Warn("Validation failed for login request")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Authenticate user
	token, user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		h.log.WithError(err).WithField("email", req.Email).Warn("Login failed")
		utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// Return token and user info
	response := map[string]interface{}{
		"access_token": token,
		"user":         user,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Login successful", response)
}

// Logout handles POST /auth/logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.userService.Logout(r.Context(), userID); err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to logout user")
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Logout failed", nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Logout successful", nil)
}

// Profile handles GET /auth/profile
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to get user profile")
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user)
}
