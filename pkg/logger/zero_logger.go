package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var zeroSinLogger *zerolog.Logger

type zeroLogger struct {
	cfg    Config
	logger *zerolog.Logger
}

var zeroLogLevelMapping = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

func newZeroLogger(cfg Config) *zeroLogger {
	logger := &zeroLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (l *zeroLogger) getLogLevel() zerolog.Level {
	level, exists := zeroLogLevelMapping[l.cfg.Level]
	if !exists {
		return zerolog.DebugLevel
	}
	return level
}

func (l *zeroLogger) Init() {
	once.Do(func() {
		loc, err := time.LoadLocation("Asia/Tehran")
		if err != nil {
			loc = time.Local
		}

		timeStamp := time.Now().In(loc).Format("2006-01-02-15-04-05")
		fileName := fmt.Sprintf("%s%s-%s.log", l.cfg.FilePath, timeStamp, uuid.New().String())

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			panic("could not open log file")
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		mw := zerolog.MultiLevelWriter(os.Stdout, file)

		var logger = zerolog.New(mw).
			With().
			Str("AppName", "2pi").
			Str("LoggerName", "Zero-log").
			Logger()

		zerolog.SetGlobalLevel(l.getLogLevel())
		zeroSinLogger = &logger
	})
	l.logger = zeroSinLogger
}

func (l *zeroLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.
		Debug().
		Time("time", time.Now()).
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.
		Debug().
		Time("time", time.Now()).
		Msgf(template, args...)
}

func (l *zeroLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.
		Info().
		Time("time", time.Now()).
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.
		Info().
		Time("time", time.Now()).
		Msgf(template, args...)
}

func (l *zeroLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.
		Warn().
		Time("time", time.Now()).
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.
		Warn().
		Time("time", time.Now()).
		Msgf(template, args...)
}

func (l *zeroLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.
		Error().
		Time("time", time.Now()).
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.
		Error().
		Time("time", time.Now()).
		Msgf(template, args...)
}

func (l *zeroLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.
		Fatal().
		Time("time", time.Now()).
		Str("Category", string(cat)).
		Str("SubCategory", string(sub)).
		Fields(logParamsToZeroParams(extra)).
		Msg(msg)
}

func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.
		Fatal().
		Time("time", time.Now()).
		Msgf(template, args...)
}

func logParamsToZeroParams(keys map[ExtraKey]interface{}) map[string]interface{} {
	params := map[string]interface{}{}

	for k, v := range keys {
		params[string(k)] = v
	}

	return params
}
