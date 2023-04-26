package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"telegram-bot/solte.lab/pkg/config"
)

func NewLogger(config *config.Logging) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		return nil, err
	}

	cw := zapcore.Lock(os.Stdout)
	je := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "log_name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	zapCore := zapcore.NewCore(je, cw, level)

	zapCore = zapcore.NewSamplerWithOptions(zapCore, time.Second, 100, 100)

	logger := zap.New(
		zapCore,
		zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel),
	)

	return logger, nil
}
