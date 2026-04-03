package cores

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func NewLogger() {
	logDir := filepath.Dir(Config().Log.Path)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   Config().Log.Path,
		MaxSize:    Config().Log.MaxSize,
		MaxBackups: Config().Log.MaxBackups,
		MaxAge:     Config().Log.MaxAge,
		Compress:   Config().Log.Compress,
	}

	var encoderConfig zapcore.EncoderConfig
	if Config().App.Env == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var consoleEncoder zapcore.Encoder
	var fileEncoder zapcore.Encoder

	if Config().App.Env == "production" {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	level, err := zapcore.ParseLevel(Config().Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	fileLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.WarnLevel
	})

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), fileLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zap.ReplaceGlobals(Logger)
}
