package handler

import (
	"context"

	"github.com/MatinHAB05/2pi/internal/domain/tokencontext"
	"github.com/MatinHAB05/2pi/pkg/logger"
)

type CommonHandler struct {
	logger logger.Logger
}

func NewCommonHandler(
	logger logger.Logger,
) CommonHandler {
	return CommonHandler{
		logger: logger,
	}
}

func (h *CommonHandler) getAuthToken(ctx context.Context) (*tokencontext.AuthenticationContextToken, bool) {
	authToken, err := tokencontext.GetTokenFromContext(ctx)
	if err != nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get token from context", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, false
	}
	return authToken, true
}

func (h *CommonHandler) getAuthTokenWithAccount(ctx context.Context) (*tokencontext.AuthenticationContextToken, bool) {
	authToken, ok := h.getAuthToken(ctx)
	if !ok {
		return nil, false
	}

	if authToken.AccountID == nil {
		h.logger.Error(logger.Handler, logger.Telegram, "failed to get account id from context - nil value", map[logger.ExtraKey]interface{}{})
		return nil, false
	}

	return authToken, true
}
