package service_contract

import (
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
)

type AuthenticationContextToken struct {
	UserId         int64
	AccountID      *int64
	AccountOwnerID *int64
}

type TokenContext struct {
	Authentication *AuthenticationContextToken
}

func MapTokenContextToService(token *tokencontext.AuthenticationContextToken) TokenContext {
	if token == nil {
		return TokenContext{}
	}

	return TokenContext{
		Authentication: &AuthenticationContextToken{
			UserId:         token.UserId,
			AccountID:      token.AccountID,
			AccountOwnerID: token.AccountOwnerID,
		},
	}
}
