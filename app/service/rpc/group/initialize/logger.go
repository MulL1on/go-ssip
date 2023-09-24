package initialize

import (
	g "go-ssip/app/service/rpc/group/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func InitLogger(serviceName string) {
	dynamicalLevel := zap.DebugLevel

	encoder := getEncoder()
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, dynamicalLevel),
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   "tmp/" + serviceName + ".log",
				MaxSize:    5,
				MaxAge:     400,
				MaxBackups: 1000,
				LocalTime:  true,
				Compress:   true,
			}),
			dynamicalLevel),
	}
	g.Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(g.Logger)
	g.Logger.Info("init logger successfully")
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stackTrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	})
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}
