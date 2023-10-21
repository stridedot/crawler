package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

// Plugin 是 zapcore.Core 的别名，两者完全相等
type Plugin = zapcore.Core

// NewLogger 用于创建一个新的 Logger
func NewLogger(plugin zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(plugin, append(DefaultOption(), options...)...)
}

func NewPlugin(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) Plugin {
	return zapcore.NewCore(DefaultEncoder(), writer, enabler)
}

// NewFilePlugin
// Lumberjack logger虽然持有File但没有暴露sync方法，
// 所以没办法利用zap的sync特性
// 额外返回一个closer，
// 保证在进程退出前close以保证写入的内容可以全部刷到到磁盘
func NewFilePlugin(filepath string, enabler zapcore.LevelEnabler) (Plugin, io.Closer) {
	writer := DefaultLumberjackLogger()
	writer.Filename = filepath
	return NewPlugin(zapcore.AddSync(writer), enabler), writer
}
