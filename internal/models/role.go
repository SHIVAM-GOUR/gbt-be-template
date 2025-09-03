package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the system
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
}

// Permission represents a permission in the system
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:100"`
	Description string         `json:"description" gorm:"size:255"`
	Resource    string         `json:"resource" gorm:"not null;size:100"` // e.g., "users", "posts"
	Action      string         `json:"action" gorm:"not null;size:50"`    // e.g., "create", "read", "update", "delete"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Roles []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint      `json:"role_id" gorm:"primaryKey"`
	PermissionID uint      `json:"permission_id" gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Role       Role       `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID"`
}

// TableName specifies the table name for the Role model
func (Role) TableName() string {
	return "roles"
}

// TableName specifies the table name for the Permission model
func (Permission) TableName() string {
	return "permissions"
}

// TableName specifies the table name for the UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}

// TableName specifies the table name for the RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}

// RoleCreateRequest represents the request payload for creating a role
type RoleCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
}

// RoleUpdateRequest represents the request payload for updating a role
type RoleUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// PermissionCreateRequest represents the request payload for creating a permission
type PermissionCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Resource    string `json:"resource" validate:"required,min=1,max=100"`
	Action      string `json:"action" validate:"required,min=1,max=50"`
}

// PermissionUpdateRequest represents the request payload for updating a permission
type PermissionUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	Resource    *string `json:"resource,omitempty" validate:"omitempty,min=1,max=100"`
	Action      *string `json:"action,omitempty" validate:"omitempty,min=1,max=50"`
}

// AssignRoleRequest represents the request payload for assigning roles to users
type AssignRoleRequest struct {
	UserID  uint   `json:"user_id" validate:"required"`
	RoleIDs []uint `json:"role_ids" validate:"required,min=1"`
}

// AssignPermissionRequest represents the request payload for assigning permissions to roles
type AssignPermissionRequest struct {
	RoleID        uint   `json:"role_id" validate:"required"`
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}

// RoleResponse represents the response payload for role data
type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

// PermissionResponse represents the response payload for permission data
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts Role model to RoleResponse
func (r *Role) ToResponse() *RoleResponse {
	resp := &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	// Convert permissions if loaded
	if len(r.Permissions) > 0 {
		resp.Permissions = make([]PermissionResponse, len(r.Permissions))
		for i, perm := range r.Permissions {
			resp.Permissions[i] = *perm.ToResponse()
		}
	}

	return resp
}

// ToResponse converts Permission model to PermissionResponse
func (p *Permission) ToResponse() *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Resource:    p.Resource,
		Action:      p.Action,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// Common permission constants
const (
	// User permissions
	PermissionUserCreate = "user.create"
	PermissionUserRead   = "user.read"
	PermissionUserUpdate = "user.update"
	PermissionUserDelete = "user.delete"

	// Role permissions
	PermissionRoleCreate = "role.create"
	PermissionRoleRead   = "role.read"
	PermissionRoleUpdate = "role.update"
	PermissionRoleDelete = "role.delete"

	// Permission permissions
	PermissionPermissionCreate = "permission.create"
	PermissionPermissionRead   = "permission.read"
	PermissionPermissionUpdate = "permission.update"
	PermissionPermissionDelete = "permission.delete"
)

// Common role constants
const (
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
	RoleUser      = "user"
)
