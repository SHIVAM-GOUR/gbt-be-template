-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create user_roles junction table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Add unique constraints
ALTER TABLE roles ADD CONSTRAINT uni_roles_name UNIQUE (name);
ALTER TABLE permissions ADD CONSTRAINT uni_permissions_name UNIQUE (name);

-- Create indexes for better performance
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_is_active ON roles(is_active);
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);
CREATE INDEX idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Insert default roles
INSERT INTO roles (name, description, is_active) VALUES
('admin', 'System administrator with full access', true),
('moderator', 'Moderator with limited administrative access', true),
('user', 'Regular user with basic access', true);

-- Insert default permissions
INSERT INTO permissions (name, description, resource, action) VALUES
-- User permissions
('user.create', 'Create new users', 'user', 'create'),
('user.read', 'View user information', 'user', 'read'),
('user.update', 'Update user information', 'user', 'update'),
('user.delete', 'Delete users', 'user', 'delete'),
('user.list', 'List all users', 'user', 'list'),

-- Role permissions
('role.create', 'Create new roles', 'role', 'create'),
('role.read', 'View role information', 'role', 'read'),
('role.update', 'Update role information', 'role', 'update'),
('role.delete', 'Delete roles', 'role', 'delete'),
('role.list', 'List all roles', 'role', 'list'),

-- Permission permissions
('permission.create', 'Create new permissions', 'permission', 'create'),
('permission.read', 'View permission information', 'permission', 'read'),
('permission.update', 'Update permission information', 'permission', 'update'),
('permission.delete', 'Delete permissions', 'permission', 'delete'),
('permission.list', 'List all permissions', 'permission', 'list');

-- Assign permissions to admin role (admin gets all permissions)
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'admin'),
    id
FROM permissions;

-- Assign limited permissions to moderator role
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'moderator'),
    id
FROM permissions
WHERE name IN ('user.read', 'user.update', 'user.list', 'role.read', 'role.list', 'permission.read', 'permission.list');

-- Assign basic permissions to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'user'),
    id
FROM permissions
WHERE name IN ('user.read');