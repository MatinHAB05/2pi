package service_contract

import (
	"context"
	"time"
)

type ShareAccountOTP struct {
	Role          string
	BaseAccountID int64
	BaseUserID    int64
}

type ShareAccountOTPService interface {
	SetShareAccountOTP(
		ctx context.Context,
		tokenContext TokenContext,
		otp string,
		data ShareAccountOTP,
		ttl time.Duration,
	) error

	GetShareAccountOTP(
		ctx context.Context,
		tokenContext TokenContext,
		otp string,
	) (*ShareAccountOTP, error)

	InvalidateShareAccountOTP(
		ctx context.Context,
		tokenContext TokenContext,
		otp string,
	) error
}
