package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DefaultEncoder 默认使用 JSON
func DefaultEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(DefaultEncoderConfig())
}

// DefaultEncoderConfig 默认的编码配置
func DefaultEncoderConfig() zapcore.EncoderConfig {
	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderConfig
}

// DefaultOption 默认的一些配置
func DefaultOption() []zap.Option {
	var stackTraceLevel zap.LevelEnablerFunc = func(level zapcore.Level) bool {
		return level >= zapcore.DPanicLevel
	}
	return []zap.Option{
		zap.AddCaller(),
		// 配置 zap 在满足特定日志级别条件时添加堆栈跟踪信息。
		// 这里，仅当日志级别大于或等于 DPanicLevel 时，才会添加堆栈跟踪
		zap.AddStacktrace(stackTraceLevel),
	}
}

// DefaultLumberjackLogger 默认的日志文件配置
func DefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:    20,
		MaxBackups: 5,
		LocalTime:  true,
		Compress:   true,
	}
}
