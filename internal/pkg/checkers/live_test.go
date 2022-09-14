package checkers_test

import (
	"testing"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/checkers"
	"github.com/stretchr/testify/assert"
)

func TestLiveChecker_Check(t *testing.T) {
	assert.NoError(t, checkers.NewLiveChecker().Check())
}
