package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MatinHAB05/2pi/config"
	service_contract "github.com/MatinHAB05/2pi/internal/application/contract"
	"github.com/MatinHAB05/2pi/pkg/logger"
	mail "github.com/MatinHAB05/2pi/pkg/wneessen-go-mail"
	"github.com/MatinHAB05/2pi/static"
)

type emailService struct {
	sender          mail.Sender
	debugModeConfig *config.DebugModeOptions
	statics         *static.StaticFiles
	appLogger       logger.Logger
}

func NewEmailService(
	sender mail.Sender,
	debugModeConfig *config.DebugModeOptions,
	statics *static.StaticFiles,
	appLogger logger.Logger,
) service_contract.EmailService {
	return &emailService{
		sender:          sender,
		debugModeConfig: debugModeConfig,
		statics:         statics,
		appLogger:       appLogger,
	}
}

func (s *emailService) SendEmail(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.SendEmailRequest) (*bool, error) {
	success := true

	if s.debugModeConfig != nil && s.debugModeConfig.Flag {
		debugText := fmt.Sprintf("[%s] To: %s | Subject: %s\nText: %s\n-------------------------------------------------------------------------------\n",
			time.Now().Format(time.RFC3339), req.To, req.Subject, req.TextBody)

		emailFile, err := os.OpenFile("./SEND_EMAIL.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			s.appLogger.Error(logger.IO, logger.Startup, "failed to open email debug file", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to open email debug log file: %w", err)
		}
		defer emailFile.Close()

		if _, err := emailFile.WriteString(debugText); err != nil {
			s.appLogger.Error(logger.IO, logger.Startup, "failed to write to email debug file", map[logger.ExtraKey]interface{}{
				logger.ErrorMessage: err.Error(),
			})
			return nil, fmt.Errorf("failed to write to email debug log file: %w", err)
		}

		return &success, nil
	}

	err := s.sender.Send(
		ctx,
		mail.SendOptions{
			To:          req.To,
			Subject:     req.Subject,
			PlainBody:   req.TextBody,
			HTMLBody:    req.HTMLBody,
			EmbedFiles:  req.EmbedFiles,
			Attachments: req.Attachments,
		},
	)
	if err != nil {
		s.appLogger.Error(logger.General, logger.ExternalService, "failed to send email via SMTP provider", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return &success, nil
}

// TODO
func (s *emailService) SendWakeUpEmail(ctx context.Context, tokenContext service_contract.TokenContext, req service_contract.SendWakeUpEmailRequest) (*bool, error) {
	buffer, err := s.statics.ExecuteWakeUpTemplate()
	if err != nil {
		s.appLogger.Error(logger.Internal, logger.ExternalService, "failed to execute WakeUp HTML template", map[logger.ExtraKey]interface{}{
			logger.ErrorMessage: err.Error(),
		})
		return nil, fmt.Errorf("failed to render  WakeUp email template: %w", err)
	}

	plainTextBody := fmt.Sprintf("...%s...", "👑MatinHAB05👑")

	success, err := s.SendEmail(ctx, tokenContext, service_contract.SendEmailRequest{
		To:       req.To,
		Subject:  "",
		HTMLBody: buffer.String(),
		TextBody: plainTextBody,
	})
	if err != nil {
		return nil, err
	}

	return success, nil
}
