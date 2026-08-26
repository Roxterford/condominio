package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
)

var loggerInstance zerolog.Logger

func Init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.CallerSkipFrameCount = 3

	var output zerolog.ConsoleWriter
	if envirotment.GetAppEnv() == envirotment.Dev {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
	}

	loggerInstance = zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	log.Logger = loggerInstance
}

func getLogger(ctx context.Context) zerolog.Logger {
	l := loggerInstance.With()
	if cid, ok := coreContext.CorrelationIDFromContext(ctx); ok && cid != "" {
		l = l.Str("correlation_id", cid)
	}
	if aid, ok := coreContext.AggregateIDFromContext(ctx); ok && aid != "" {
		l = l.Str("aggregate_id", aid)
	}
	return l.Logger()
}

func DebugCtx(ctx context.Context, msg string, fields ...any) {
	l := getLogger(ctx)
	l.Debug().Caller().Fields(fieldsToMap(fields)).Msg(msg)
}

func InfoCtx(ctx context.Context, msg string, fields ...any) {
	l := getLogger(ctx)
	l.Info().Caller().Fields(fieldsToMap(fields)).Msg(msg)
}

func WarnCtx(ctx context.Context, msg string, fields ...any) {
	l := getLogger(ctx)
	l.Warn().Caller().Fields(fieldsToMap(fields)).Msg(msg)
}

func ErrorCtx(ctx context.Context, err error, msg string, fields ...any) {
	l := getLogger(ctx)
	l.Error().Caller().Err(err).Fields(fieldsToMap(fields)).Msg(msg)
}

func FatalCtx(ctx context.Context, err error, msg string, fields ...any) {
	l := getLogger(ctx)
	l.Fatal().Caller().Err(err).Fields(fieldsToMap(fields)).Msg(msg)
}

func With(ctx context.Context, fields ...any) zerolog.Logger {
	return getLogger(ctx).With().Fields(fieldsToMap(fields)).Logger()
}

func fieldsToMap(fields []any) map[string]any {
	if len(fields)%2 != 0 {
		return map[string]any{"invalid_fields": fields}
	}
	m := make(map[string]any, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		if key, ok := fields[i].(string); ok {
			m[key] = fields[i+1]
		}
	}
	return m
}

// For backward compatibility with existing logger.Error calls
func LegacyError(err error) {
	loggerInstance.Error().Err(err).Msg("legacy error")
}

func LegacyWarning(msg string, args ...any) {
	loggerInstance.Warn().Msgf(msg, args...)
}

func LegacyInfo(msg string, args ...any) {
	loggerInstance.Info().Msgf(msg, args...)
}

// Backward compatible printf-style functions (without context)
func PrintfDebug(msg string, args ...any) {
	loggerInstance.Debug().Msgf(msg, args...)
}

func PrintfInfo(msg string, args ...any) {
	loggerInstance.Info().Msgf(msg, args...)
}

func PrintfWarn(msg string, args ...any) {
	loggerInstance.Warn().Msgf(msg, args...)
}

func PrintfError(msg string, args ...any) {
	loggerInstance.Error().Msgf(msg, args...)
}

// Deprecated: use Debug(ctx, ...) instead
func Debugf(msg string, args ...any) {
	loggerInstance.Debug().Msgf(msg, args...)
}

// Deprecated: use Info(ctx, ...) instead
func Infof(msg string, args ...any) {
	loggerInstance.Info().Msgf(msg, args...)
}

// Deprecated: use Warn(ctx, ...) instead
func Warnf(msg string, args ...any) {
	loggerInstance.Warn().Msgf(msg, args...)
}

// Deprecated: use Error(ctx, ...) instead
func Errorf(msg string, args ...any) {
	loggerInstance.Error().Msgf(msg, args...)
}

// Backward compatible printf-style functions (matching old API exactly)
func Debug(msg string, args ...any) {
	loggerInstance.Debug().Msgf(msg, args...)
}

func Info(msg string, args ...any) {
	loggerInstance.Info().Msgf(msg, args...)
}

func Warning(msg string, args ...any) {
	loggerInstance.Warn().Msgf(msg, args...)
}

func Error(msg string, args ...any) {
	loggerInstance.Error().Msgf(msg, args...)
}

func Fatal(msg string, args ...any) {
	loggerInstance.Fatal().Msgf(msg, args...)
}

func Fatalf(msg string, args ...any) {
	loggerInstance.Fatal().Msgf(msg, args...)
}
