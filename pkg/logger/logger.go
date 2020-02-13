package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

// default info level
var level = zap.NewAtomicLevelAt(zap.InfoLevel)

func init() {
	config := zap.Config{
		Level:             level,
		Development:       false,
		DisableStacktrace: true,
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stdout"}, //"outputPaths": ["stdout", "/tmp/logs"],
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %v", err))
	}
	log = logger.Sugar()
}

//SetDebug ...
func SetDebug() {
	level.SetLevel(zapcore.DebugLevel)
}

//SetLevel ...
func SetLevel(lvl zapcore.Level) {
	level.SetLevel(lvl)
}

//IsDebug ...
func IsDebug() bool {
	return level.Level() == zapcore.DebugLevel
}

//GetLogger ...
func GetLogger() *zap.SugaredLogger {
	return log
}
