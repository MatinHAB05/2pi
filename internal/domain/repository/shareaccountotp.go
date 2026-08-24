package repository_contract

import (
	"context"
	"time"
)

type ShareAccountOTP struct {
	Role          string
	BaseAccountID string
	BaseUserID    string
}

type ShareAccountOTPRepository interface {
	Set(ctx context.Context, otp string, cache *ShareAccountOTP, ttl time.Duration) error
	Get(ctx context.Context, otp string) (*ShareAccountOTP, error)
	Delete(ctx context.Context, otp string) (bool, error)
	Exists(ctx context.Context, otp string) (*bool, error)
}
