package exception

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrTargetAccountNotFound = errors.New("user-account not found")
	ErrNilEnforcer           = errors.New("RBAC enforcer or underlying Casbin instance is nil")
	ErrEmptyUserID           = errors.New("user ID cannot be empty")
	ErrEmptyTargetAccount    = errors.New("target account ID cannot be empty")
	ErrEmptyRole             = errors.New("role cannot be empty")
	ErrEmptyAction           = errors.New("action cannot be empty")
)
