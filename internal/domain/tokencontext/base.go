package tokencontext

import "context"

type contextKey string

type TokenContext struct {
	Authentication *AuthenticationContextToken
	Info           *UserInfoContextToken
}

func GetTokenContextFromContext(ctx context.Context) (*TokenContext, error) {
	auth, err := GetAuthenticationTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	info, err := GetInfoTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &TokenContext{
		Authentication: auth,
		Info:           info,
	}, nil
}

func SetTokenContextInContext(ctx context.Context, tokenContext *TokenContext) context.Context {
	if tokenContext == nil {
		return ctx
	}

	if tokenContext.Authentication != nil {
		ctx = SetAuthenticationTokenInContext(ctx, tokenContext.Authentication)
	}

	if tokenContext.Info != nil {
		ctx = SetInfoTokenInContext(ctx, tokenContext.Info)
	}

	return ctx
}
