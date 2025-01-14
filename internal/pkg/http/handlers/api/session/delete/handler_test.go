package delete_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api/session/delete"

	"github.com/gorilla/mux"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	var cases = []struct {
		name        string
		giveReqVars func(sessionUUID string) map[string]string
		checkResult func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "without registered session UUID",
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
				assert.JSONEq(t,
					`{"code":500,"success":false,"message":"cannot extract session UUID"}`, rr.Body.String(),
				)
			},
		},
		{
			name: "session was not found",
			giveReqVars: func(_ string) map[string]string {
				return map[string]string{"sessionUUID": "aa-bb-cc-dd"}
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.JSONEq(t,
					`{"code":404,"success":false,"message":"session with UUID aa-bb-cc-dd was not found"}`, rr.Body.String(),
				)
			},
		},
		{
			name: "success",
			giveReqVars: func(sessionUUID string) map[string]string {
				return map[string]string{"sessionUUID": sessionUUID}
			},
			checkResult: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.JSONEq(t,
					`{"success":true}`, rr.Body.String(),
				)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewInMemory(time.Minute, 10)
			defer s.Close()

			sessionUUID, err := s.CreateSession([]byte{}, 201, "", 0)
			assert.NoError(t, err)

			var (
				req, _  = http.NewRequest(http.MethodPost, "http://testing", http.NoBody)
				rr      = httptest.NewRecorder()
				handler = delete.NewHandler(s)
			)

			if tt.giveReqVars != nil {
				req = mux.SetURLVars(req, tt.giveReqVars(sessionUUID))
			}

			handler.ServeHTTP(rr, req)

			tt.checkResult(t, rr)
		})
	}
}
