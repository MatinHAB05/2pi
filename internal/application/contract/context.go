package service_contract

import (
	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
)

type AuthenticationContextToken struct {
	UserId         int64
	AccountID      *int64
	AccountOwnerID *int64
}

type UserInfoContextToken struct {
	Lang entity.Lang
}

type TokenContext struct {
	Authentication *AuthenticationContextToken
	Info           *UserInfoContextToken
}

func MapTokenContextToService(tokenCtx *tokencontext.TokenContext) TokenContext {
	if tokenCtx == nil {
		return TokenContext{}
	}

	var auth *AuthenticationContextToken
	if tokenCtx.Authentication != nil {
		auth = &AuthenticationContextToken{
			UserId:         tokenCtx.Authentication.UserId,
			AccountID:      tokenCtx.Authentication.AccountID,
			AccountOwnerID: tokenCtx.Authentication.AccountOwnerID,
		}
	}

	var info *UserInfoContextToken
	if tokenCtx.Info != nil {
		info = &UserInfoContextToken{
			Lang: tokenCtx.Info.Lang,
		}
	}

	return TokenContext{
		Authentication: auth,
		Info:           info,
	}
}

// MapTokenContextToServiceJustAuth wraps an AuthenticationContextToken into service TokenContext.
func MapTokenContextToServiceJustAuth(tokenCtx *tokencontext.AuthenticationContextToken) TokenContext {
	if tokenCtx == nil {
		return TokenContext{}
	}

	return TokenContext{
		Authentication: &AuthenticationContextToken{
			UserId:         tokenCtx.UserId,
			AccountID:      tokenCtx.AccountID,
			AccountOwnerID: tokenCtx.AccountOwnerID,
		},
	}
}

// MapTokenContextToServiceJustInfo wraps an UserInfoContextToken into service TokenContext.
func MapTokenContextToServiceJustInfo(tokenCtx *tokencontext.UserInfoContextToken) TokenContext {
	if tokenCtx == nil {
		return TokenContext{}
	}

	return TokenContext{
		Info: &UserInfoContextToken{
			Lang: tokenCtx.Lang,
		},
	}
}

// MapAuthenticationTokenToService maps domain AuthenticationContextToken directly to service_contract.
func MapAuthenticationTokenToService(token *tokencontext.AuthenticationContextToken) *AuthenticationContextToken {
	if token == nil {
		return nil
	}
	return &AuthenticationContextToken{
		UserId:         token.UserId,
		AccountID:      token.AccountID,
		AccountOwnerID: token.AccountOwnerID,
	}
}

// MapInfoTokenToService maps domain UserInfoContextToken directly to service_contract.
func MapInfoTokenToService(token *tokencontext.UserInfoContextToken) *UserInfoContextToken {
	if token == nil {
		return nil
	}
	return &UserInfoContextToken{
		Lang: token.Lang,
	}
}
