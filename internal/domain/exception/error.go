package exception

import "errors"

var (
	ErrUserNotFound                      = errors.New("user not found")
	ErrShareAccountAccessOTPCodeNotFound = errors.New("share account access otp code not found")
	ErrTargetAccountNotFound             = errors.New("user-account not found")
	ErrNilEnforcer                       = errors.New("RBAC enforcer or underlying Casbin instance is nil")
	ErrEmptyUserID                       = errors.New("user ID cannot be empty")
	ErrEmptyTargetAccount                = errors.New("target account ID cannot be empty")
	ErrEmptyRole                         = errors.New("role cannot be empty")
	ErrEmptyAction                       = errors.New("action cannot be empty")
	ErrUserIDNotFound                    = errors.New("user_id missing from context")
	ErrInvalidUserIDType                 = errors.New("user_id in context has invalid type")
	ErrLangNotFound                      = errors.New("lang not found in context")
)
