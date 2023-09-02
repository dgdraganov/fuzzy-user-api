package log

import (
	"fmt"
	"strings"

	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type ILogger interface {
// 	Infow(string, ...interface{})
// 	Errorw(string, ...interface{})
// 	Warnw(string, ...interface{})
// 	Debugw(string, ...interface{})

// 	Info(...interface{})
// 	Error(...interface{})
// 	Warn(...interface{})
// 	Debug(...interface{})

// 	Infof(string, ...interface{})
// 	Errorf(string, ...interface{})
// 	Warnf(string, ...interface{})
// 	Debugf(string, ...interface{})
// }

const (
	Info  string = "info"
	Warn  string = "info"
	Debug string = "info"
)

type ZapLogger struct {
	*zap.SugaredLogger
}

func NewZapLogger(serviceName string, logLevel zapcore.Level) *ZapLogger {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(logLevel)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339Nano))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:       atom,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	core, err := config.Build()
	if err != nil {
		panic("can not build logger config")
	}

	logger := zap.New(core.Core(), zap.AddCaller())

	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Sugar().Info("logger sync completed")
		}
	}()

	sugar := logger.Sugar()

	// setting default fields
	sugar = sugar.With(
		"microservice", serviceName,
	)

	return &ZapLogger{
		sugar,
	}
}

/*
Possible values:
ERROR, WARN, INFO, DEBUG
*/
func parseLogLevel(level string) zapcore.Level {
	level = strings.ToUpper(level)
	var newLevel zapcore.Level

	switch level {
	case "ERROR":
		newLevel = zap.ErrorLevel
	case "WARN":
		newLevel = zap.WarnLevel
	case "INFO":
		newLevel = zap.InfoLevel
	case "DEBUG":
		newLevel = zap.DebugLevel
	default:
		panic(fmt.Sprintf("invalid log level %s", level))
	}
	return newLevel
}
