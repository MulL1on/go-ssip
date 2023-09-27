package initialize

import (
	g "go-ssip/app/service/rpc/user/global"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func InitLogger(serviceName string) {
	dynamicalLevel := zap.DebugLevel

	encoder := getECSEncoder()
	cores := [...]zapcore.Core{
		ecszap.NewCore(encoder, os.Stdout, dynamicalLevel),
		ecszap.NewCore(
			getECSEncoder(),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   "tmp/" + serviceName + ".json",
				MaxSize:    5,
				MaxAge:     400,
				MaxBackups: 1000,
				LocalTime:  true,
				Compress:   true,
			}),
			dynamicalLevel),
	}
	g.Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller()).With(zap.String("app", serviceName))

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

func getECSEncoder() ecszap.EncoderConfig {
	cfg := ecszap.NewDefaultEncoderConfig()
	cfg.EncodeCaller = zapcore.FullCallerEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EnableStackTrace = true
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	cfg.LineEnding = zapcore.DefaultLineEnding
	cfg.EnableCaller = true
	return cfg
}
