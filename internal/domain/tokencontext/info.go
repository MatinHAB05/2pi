package tokencontext

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/entity"
	"github.com/MatinHAB05/2pi/internal/domain/exception"
)

const (
	LangKey contextKey = "lang"
)

type UserInfoContextToken struct {
	Lang entity.Lang
}

func GetInfoTokenFromContext(ctx context.Context) (*UserInfoContextToken, error) {
	lang, ok := ctx.Value(LangKey).(entity.Lang)
	if !ok {
		return nil, exception.ErrLangNotFound
	}
	return &UserInfoContextToken{Lang: lang}, nil
}

func SetInfoTokenInContext(ctx context.Context, token *UserInfoContextToken) context.Context {
	if token == nil {
		return ctx
	}
	return context.WithValue(ctx, LangKey, token.Lang)
}
