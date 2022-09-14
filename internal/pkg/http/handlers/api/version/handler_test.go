package version_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api/version"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	var (
		req, _ = http.NewRequest(http.MethodGet, "http://testing", http.NoBody)
		rr     = httptest.NewRecorder()
	)

	version.NewHandler("1.2.3@foo")(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.JSONEq(t, `{"version":"1.2.3@foo"}`, rr.Body.String())
}
