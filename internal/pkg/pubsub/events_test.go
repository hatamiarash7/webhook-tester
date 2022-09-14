package pubsub_test

import (
	"testing"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestRegisteredEvent(t *testing.T) {
	e := pubsub.NewRequestRegisteredEvent("foo")

	assert.Equal(t, []byte("foo"), e.Data())
	assert.Equal(t, "request-registered", e.Name())
}

func TestNewRequestDeletedEvent(t *testing.T) {
	e := pubsub.NewRequestDeletedEvent("foo")

	assert.Equal(t, []byte("foo"), e.Data())
	assert.Equal(t, "request-deleted", e.Name())
}

func TestNewAllRequestsDeletedEvent(t *testing.T) {
	e := pubsub.NewAllRequestsDeletedEvent()

	assert.Equal(t, []byte("*"), e.Data())
	assert.Equal(t, "requests-deleted", e.Name())
}
