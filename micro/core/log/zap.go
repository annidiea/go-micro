package log

import (
	"fmt"
	"io"

	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(config *Config) *zap.Logger {
	fmt.Println(config)
	// 自定义zap日志配置
	encoderconfig := zap.NewProductionEncoderConfig()
	// 自定义时间格式
	//【shopstar】 2006-01-02 15:04:05
	encoderconfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format(config.Prefix + "2006-01-02 15:04:05"))
	}
	// 创建普通的编译器
	encoder := zapcore.NewConsoleEncoder(encoderconfig)
	// 创建写入器
	writer := zapcore.AddSync(getWriter(config))
	// new 创建的是zap的核心对象
	// 编译器，写入器，参数级别
	core := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	// 1. 核心（编译器，写入器，参数级别）
	return zap.New(core, zap.AddCaller(), zap.Hooks())
}

func getWriter(config *Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.Maxsize,
		MaxBackups: config.Maxbackups,
		MaxAge:     config.Maxage,
		Compress:   config.Compress,
	}
}
