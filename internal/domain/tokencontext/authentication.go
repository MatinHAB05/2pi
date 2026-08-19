package tokencontext

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/exception"
)

const (
	UserIDKey         contextKey = "user_id"
	AccountIDKey      contextKey = "account_id"
	AccountOwnerIDKey contextKey = "account_owner_id"
)

type AuthenticationContextToken struct {
	UserId         int64
	AccountID      *int64
	AccountOwnerID *int64
}

// GetAuthenticationTokenFromContext constructs an AuthenticationContextToken from context key-value pairs.
func GetAuthenticationTokenFromContext(ctx context.Context) (*AuthenticationContextToken, error) {
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

func SetAuthenticationTokenInContext(ctx context.Context, token *AuthenticationContextToken) context.Context {
	if token == nil {
		return ctx
	}

	ctx = context.WithValue(ctx, UserIDKey, token.UserId)

	if token.AccountID != nil {
		ctx = context.WithValue(ctx, AccountIDKey, token.AccountID)
	}

	if token.AccountOwnerID != nil {
		ctx = context.WithValue(ctx, AccountOwnerIDKey, token.AccountOwnerID)
	}

	return ctx
}
