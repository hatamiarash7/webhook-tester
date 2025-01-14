package settings_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/config"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api/settings"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	var cases = []struct {
		name        string
		setUp       func(cfg *config.Config)
		checkResult func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "without registered session UUID",
			setUp: func(cfg *config.Config) {
				cfg.MaxRequests = 123
				cfg.SessionTTL = time.Second * 321
				cfg.MaxRequestBodySize = 222
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.JSONEq(t, `{
					"limits": {"max_requests":123, "session_lifetime_sec":321, "max_webhook_body_size": 222}
				}`, rr.Body.String())
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				req, _ = http.NewRequest(http.MethodPost, "http://testing", http.NoBody)
				rr     = httptest.NewRecorder()
				cfg    = config.Config{}
			)

			if tt.setUp != nil {
				tt.setUp(&cfg)
			}

			handler := settings.NewGetSettingsHandler(cfg)

			handler.ServeHTTP(rr, req)

			tt.checkResult(t, rr)
		})
	}
}

func BenchmarkNewGetSettingsHandler(b *testing.B) {
	b.ReportAllocs()

	var (
		req, _ = http.NewRequest(http.MethodPost, "http://testing", http.NoBody)
		rr     = httptest.NewRecorder()
		cfg    = config.Config{}
	)

	h := settings.NewGetSettingsHandler(cfg)

	for n := 0; n < b.N; n++ {
		h.ServeHTTP(rr, req)
	}
}
