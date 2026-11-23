package cores

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger        *zap.Logger
	loggerMu      sync.Mutex
	fileLogCloser *lumberjack.Logger
)

func NewLogger() {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if Logger != nil && fileLogCloser != nil {
		return
	}
	initLoggerLocked()
}

// CloseLogger flushes buffered log entries and closes the underlying rotating file handle if open.
func CloseLogger() error {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if Logger != nil {
		_ = Logger.Sync()
	}
	if fileLogCloser != nil {
		err := fileLogCloser.Close()
		fileLogCloser = nil
		nop := zap.NewNop()
		Logger = nop
		zap.ReplaceGlobals(nop)
		return err
	}
	return nil
}

func initLoggerLocked() {
	logDir := filepath.Dir(Config().Log.Path)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0750); err != nil {
			_, _ = os.Stderr.WriteString("Logger initialization warning: failed to create log directory: " + err.Error() + "\n")
		}
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   Config().Log.Path,
		MaxSize:    Config().Log.MaxSize,
		MaxBackups: Config().Log.MaxBackups,
		MaxAge:     Config().Log.MaxAge,
		Compress:   Config().Log.Compress,
	}
	fileLogCloser = lumberJackLogger

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
		_, _ = os.Stderr.WriteString("Logger initialization warning: invalid LOG_LEVEL \"" + Config().Log.Level + "\", falling back to \"info\"\n")
		level = zapcore.InfoLevel
	}

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	fileLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= level
	})

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), fileLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	zap.ReplaceGlobals(Logger)
}
