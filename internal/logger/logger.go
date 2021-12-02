package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// layoutDateFormatLogFile - date format for file logging
const layoutDateFormatLogFile = "2006-01-02T15:04:05.999"

type ctxKey struct{}

var attachedLoggerKey = &ctxKey{}
var globalLogger *zap.SugaredLogger
var notSugaredLogger *zap.Logger

// GetNotSugaredLogger - get not suggared logger
func GetNotSugaredLogger() *zap.Logger {
	return notSugaredLogger
}

// InitLogger - create logger
func InitLogger(cfg *config.Config) (syncFn func()) {

	EncoderConfig := zap.NewProductionEncoderConfig()
	LogLevel := zap.InfoLevel

	if cfg.Project.Debug {
		EncoderConfig = zap.NewDevelopmentEncoderConfig()
		LogLevel = zap.DebugLevel
	}

	EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(layoutDateFormatLogFile))
	})

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(EncoderConfig),
		os.Stderr,
		zap.NewAtomicLevelAt(LogLevel),
	)

	notSugaredLogger = zap.New(consoleCore)

	sugaredLogger := notSugaredLogger.Sugar()
	SetLogger(sugaredLogger.With(
		"service", cfg.Project.Name,
	))

	return func() {
		if err := notSugaredLogger.Sync(); err != nil {
			ErrorKV(context.Background(), "not suggared logger sync", "err", err)
		}
	}
}

// StrToZapLevel - convert string to zap.level
func StrToZapLevel(str string) (zapcore.Level, bool) {

	switch strings.ToLower(str) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info":
		return zapcore.InfoLevel, true
	case "warn":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	default:
		return zapcore.InfoLevel, false
	}
}

func fromContext(ctx context.Context) *zap.SugaredLogger {
	var result *zap.SugaredLogger
	if attachedLogger, ok := ctx.Value(attachedLoggerKey).(*zap.SugaredLogger); ok {
		result = attachedLogger
	} else {
		result = globalLogger
	}

	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		if spanCtx, ok := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext); ok {
			result = result.With("trace-id", spanCtx.TraceID())
		}
	}

	return result
}

// ErrorKV - logger add string log level error
func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Errorw(message, kvs...)
}

// WarnKV - logger add string log level warn
func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Warnw(message, kvs...)
}

// InfoKV - logger add string log level info
func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Infow(message, kvs...)
}

// DebugKV - logger add string log level debug
func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Debugw(message, kvs...)
}

// FatalKV - logger add string log level fatal
func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	fromContext(ctx).Fatalw(message, kvs...)
}

// AttachLogger - context attach logger
func AttachLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, attachedLoggerKey, logger)
}

// CloneWithLevel - clone logger with level
func CloneWithLevel(ctx context.Context, newLevel zapcore.Level) *zap.SugaredLogger {
	return fromContext(ctx).
		Desugar().
		WithOptions(WithLevel(newLevel)).
		Sugar()
}

// SetLogger - set global logger
func SetLogger(newLogger *zap.SugaredLogger) {
	globalLogger = newLogger
}

func init() {
	notSugaredLogger, err := zap.NewProduction()
	if err != nil {
		FatalKV(context.Background(), "Fatal init logger", err)
	}

	globalLogger = notSugaredLogger.Sugar()
}
