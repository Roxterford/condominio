package logger

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupTestLogger() *bytes.Buffer {
	os.Setenv("APP_ENV", "dev")
	var buf bytes.Buffer
	Init()
	loggerInstance = loggerInstance.Output(&buf)
	log.Logger = loggerInstance
	return &buf
}

func TestInfoCtxShowsCallerLine(t *testing.T) {
	buf := setupTestLogger()
	ctx := context.Background()

	InfoCtx(ctx, "test message") // THIS LINE should appear in the log

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Fatalf("expected 'test message' in output, got: %s", output)
	}

	if !strings.Contains(output, "logger_test.go") {
		t.Fatalf("expected 'logger_test.go' in caller field, got: %s", output)
	}

	if strings.Contains(output, "event.go") {
		t.Fatalf("caller should NOT be event.go (zerolog internal), got: %s", output)
	}

	if strings.Contains(output, "logger.go") {
		t.Fatalf("caller should NOT be logger.go, got: %s", output)
	}
}

func TestErrorCtxShowsCallerLine(t *testing.T) {
	buf := setupTestLogger()
	ctx := context.Background()

	ErrorCtx(ctx, nil, "error test") // THIS LINE should appear in the log

	output := buf.String()
	if !strings.Contains(output, "error test") {
		t.Fatalf("expected 'error test' in output, got: %s", output)
	}

	if !strings.Contains(output, "logger_test.go") {
		t.Fatalf("expected 'logger_test.go' in caller field, got: %s", output)
	}

	if strings.Contains(output, "event.go") {
		t.Fatalf("caller should NOT be event.go (zerolog internal), got: %s", output)
	}
}

func TestWarnCtxShowsCallerLine(t *testing.T) {
	buf := setupTestLogger()
	ctx := context.Background()

	WarnCtx(ctx, "warn test") // THIS LINE should appear in the log

	output := buf.String()
	if !strings.Contains(output, "logger_test.go") {
		t.Fatalf("expected 'logger_test.go' in caller field, got: %s", output)
	}

	if strings.Contains(output, "event.go") {
		t.Fatalf("caller should NOT be event.go (zerolog internal), got: %s", output)
	}
}

func TestDebugCtxShowsCallerLine(t *testing.T) {
	os.Setenv("APP_ENV", "dev")
	var buf bytes.Buffer
	Init()
	loggerInstance = loggerInstance.Output(&buf).Level(zerolog.DebugLevel)
	log.Logger = loggerInstance

	ctx := context.Background()

	DebugCtx(ctx, "debug test") // THIS LINE should appear in the log

	output := buf.String()
	if !strings.Contains(output, "logger_test.go") {
		t.Fatalf("expected 'logger_test.go' in caller field, got: %s", output)
	}

	if strings.Contains(output, "event.go") {
		t.Fatalf("caller should NOT be event.go (zerolog internal), got: %s", output)
	}
}
