package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New zap sugared logger
func New(setters ...Option) *zap.SugaredLogger {
	// Default Options
	args := &Options{
		Level:   zap.ErrorLevel,
		Encoder: zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	}
	for _, setter := range setters {
		setter(args)
	}

	loggingLevel := zap.NewAtomicLevel()
	loggingLevel.SetLevel(args.Level)

	return zap.New(zapcore.NewCore(
		args.Encoder,
		zapcore.Lock(os.Stdout),
		loggingLevel,
	)).Sugar()
}

type Options struct {
	Level   zapcore.Level
	Encoder zapcore.Encoder
}

type Option func(*Options)

func Level(level string) Option {
	return func(args *Options) {
		var lv zapcore.Level
		switch level {
		case "panic":
			lv = zap.PanicLevel
		case "fatal":
			lv = zap.FatalLevel
		case "warn":
			lv = zap.WarnLevel
		case "info":
			lv = zap.InfoLevel
		case "debug":
			lv = zap.DebugLevel
		default:
			lv = zap.ErrorLevel
		}
		args.Level = lv
	}
}

func Encoder(format string) Option {
	return func(args *Options) {
		var enc zapcore.Encoder
		switch format {
		case "console":
			enc = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		default:
			enc = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}
		args.Encoder = enc
	}
}
