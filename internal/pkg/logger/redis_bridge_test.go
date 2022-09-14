package logger_test

import (
	"context"
	"testing"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/logger"
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
)

func TestRedisBridge_Printf(t *testing.T) {
	output := capturer.CaptureStderr(func() {
		log, err := logger.New(false, false, false)
		assert.NoError(t, err)

		br := logger.NewRedisBridge(log)

		br.Printf(context.Background(), "%s", "foobar")
	})

	assert.Contains(t, output, "warn")
	assert.Contains(t, output, "foobar")
}
