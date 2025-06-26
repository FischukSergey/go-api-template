package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger инициализирует глобальный логгер с заданным уровнем логирования.
func InitLogger(level string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		return fmt.Errorf("unsupported log level: %s", level)
	}

	var config zap.Config

	if level == "debug" {
		// Для debug - JSON формат с подробной информацией
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel)
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // строчные цветные
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		config.DisableStacktrace = false
		config.DisableCaller = false
	} else {
		// Для info, warn, error - текстовый формат
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel)
		config.Encoding = "json"
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableStacktrace = true // Отключаем стектрейсы для обычных логов
		config.DisableCaller = true     // Отключаем caller для чистоты вывода
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	// Устанавливаем созданный логгер как глобальный для zap.L()
	zap.ReplaceGlobals(Logger)

	return nil
}

// GetLogger возвращает глобальный логгер.
func GetLogger() *zap.Logger {
	if Logger == nil {
		// Fallback на случай если InitLogger не был вызван
		Logger, _ = zap.NewProduction()
	}
	return Logger
}

// Sync синхронизирует логгер (важно вызывать при завершении приложения).
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
