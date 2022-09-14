package metrics_test

import (
	"testing"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWebSockets_Register(t *testing.T) {
	var (
		registry = prometheus.NewRegistry()
		ws       = metrics.NewWebsockets()
	)

	assert.NoError(t, ws.Register(registry))

	count, err := testutil.GatherAndCount(registry, "websockets_active_clients_count")
	assert.NoError(t, err)

	assert.Equal(t, 1, count)
}

func TestWebSockets_IncrementActiveClients(t *testing.T) {
	ws := metrics.NewWebsockets()

	ws.IncrementActiveClients()

	metric := getMetric(&ws, "websockets_active_clients_count")
	assert.Equal(t, float64(1), metric.Gauge.GetValue())
}

func TestWebSockets_DecrementActiveClients(t *testing.T) {
	ws := metrics.NewWebsockets()

	ws.DecrementActiveClients()

	metric := getMetric(&ws, "websockets_active_clients_count")
	assert.Equal(t, float64(-1), metric.Gauge.GetValue())
}
