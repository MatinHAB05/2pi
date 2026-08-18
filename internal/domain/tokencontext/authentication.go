package tokencontext

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/exception"
)

const (
	UserIDKey         contextKey = "user_id"
	CompletedKey      contextKey = "completed"
	AccountIDKey      contextKey = "account_id"
	AccountOwnerIDKey contextKey = "account_owner_id"
)

type AuthenticationContextToken struct {
	UserId         int64
	Completed      bool
	AccountID      *int64
	AccountOwnerID *int64
}

// GetTokenFromContext constructs an AuthenticationContextToken from context key-value pairs.
func GetTokenFromContext(ctx context.Context) (*AuthenticationContextToken, error) {
	rawUserID := ctx.Value(UserIDKey)
	if rawUserID == nil {
		return nil, exception.ErrUserIDNotFound
	}

	var userID int64
	switch v := rawUserID.(type) {
	case int64:
		userID = v
	case *int64:
		if v == nil {
			return nil, exception.ErrUserIDNotFound
		}
		userID = *v
	case int:
		userID = int64(v)
	default:
		return nil, exception.ErrInvalidUserIDType
	}

	token := &AuthenticationContextToken{
		UserId: userID,
	}

	rawCompleted := ctx.Value(CompletedKey)
	if rawCompleted == nil {
		return nil, exception.ErrCompletedNotFound
	}

	var completed bool
	switch v := rawCompleted.(type) {
	case bool:
		completed = v
	default:
		return nil, exception.ErrInvalidCompletedType
	}
	token.Completed = completed

	// Extract optional AccountID
	if rawAccID := ctx.Value(AccountIDKey); rawAccID != nil {
		switch v := rawAccID.(type) {
		case *int64:
			token.AccountID = v
		case int64:
			token.AccountID = &v
		}
	}

	// Extract optional AccountOwnerID
	if rawOwnerID := ctx.Value(AccountOwnerIDKey); rawOwnerID != nil {
		switch v := rawOwnerID.(type) {
		case *int64:
			token.AccountOwnerID = v
		case int64:
			token.AccountOwnerID = &v
		}
	}

	return token, nil
}

func SetTokenInContext(ctx context.Context, token *AuthenticationContextToken) context.Context {
	if token == nil {
		return ctx
	}

	ctx = context.WithValue(ctx, UserIDKey, token.UserId)
	ctx = context.WithValue(ctx, CompletedKey, token.Completed)

	if token.AccountID != nil {
		ctx = context.WithValue(ctx, AccountIDKey, token.AccountID)
	}

	if token.AccountOwnerID != nil {
		ctx = context.WithValue(ctx, AccountOwnerIDKey, token.AccountOwnerID)
	}

	return ctx
}
