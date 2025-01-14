package all_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api/session/requests/all"

	"github.com/gorilla/mux"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTPRequestErrors(t *testing.T) {
	var cases = []struct {
		name           string
		giveReqVars    map[string]string
		wantStatusCode int
		wantJSON       string
	}{
		{
			name:           "without registered session UUID",
			giveReqVars:    nil,
			wantStatusCode: http.StatusInternalServerError,
			wantJSON:       `{"code":500,"success":false,"message":"cannot extract session UUID"}`,
		},
		{
			name:           "session was not found",
			giveReqVars:    map[string]string{"sessionUUID": "aa-bb-cc-dd"},
			wantStatusCode: http.StatusNotFound,
			wantJSON:       `{"code":404,"success":false,"message":"session with UUID aa-bb-cc-dd was not found"}`,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewInMemory(time.Minute, 1)
			defer s.Close()

			var (
				req, _  = http.NewRequest(http.MethodPost, "http://testing", http.NoBody)
				rr      = httptest.NewRecorder()
				handler = all.NewHandler(s)
			)

			if tt.giveReqVars != nil {
				req = mux.SetURLVars(req, tt.giveReqVars)
			}

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			assert.JSONEq(t, tt.wantJSON, rr.Body.String())
		})
	}
}

func TestHandler_ServeHTTPSuccessSingle(t *testing.T) {
	s := storage.NewInMemory(time.Minute, 10)
	defer s.Close()

	var (
		req, _  = http.NewRequest(http.MethodGet, "http://test", http.NoBody)
		rr      = httptest.NewRecorder()
		handler = all.NewHandler(s)
	)

	// create session
	sessionUUID, err := s.CreateSession([]byte("foo"), 202, "foo/bar", 0)
	assert.NoError(t, err)

	// create ONE request for the session
	requestUUID, err := s.CreateRequest(
		sessionUUID,
		"1.2.2.1",
		"PUT",
		"http://example.com/foo",
		[]byte("foobar"),
		map[string]string{"aaa": "bar", "bbb": "foo"},
	)
	assert.NoError(t, err)

	request, _ := s.GetRequest(sessionUUID, requestUUID)

	req = mux.SetURLVars(req, map[string]string{"sessionUUID": sessionUUID})

	handler.ServeHTTP(rr, req)

	runtime.Gosched()
	<-time.After(time.Millisecond) // goroutine must be done

	assert.JSONEq(t, `[{
		"client_address":"1.2.2.1",
		"content_base64":"`+base64.StdEncoding.EncodeToString([]byte("foobar"))+`",
		"created_at_unix":`+strconv.FormatInt(request.CreatedAt().Unix(), 10)+`,
		"headers":[{"name": "aaa", "value": "bar"},{"name": "bbb", "value": "foo"}],
		"method":"PUT",
		"url":"http://example.com/foo",
		"uuid":"`+request.UUID()+`"
	}]`, rr.Body.String())
}

func TestHandler_ServeHTTPSuccessMultiple(t *testing.T) { // must be sorted
	s := storage.NewInMemory(time.Minute, 3)
	defer s.Close()

	var (
		req, _  = http.NewRequest(http.MethodGet, "http://test", http.NoBody)
		rr      = httptest.NewRecorder()
		handler = all.NewHandler(s)
	)

	// create session
	sessionUUID, err := s.CreateSession([]byte("foo"), 202, "foo/bar", 0)
	assert.NoError(t, err)

	// create THREE requests for the session
	_, _ = s.CreateRequest( // must be ignored, storage limit = 3
		sessionUUID,
		"1.1.1.1",
		"PUT",
		"http://example.com/foo1",
		[]byte("foobar"),
		nil,
	)
	request1UUID, _ := s.CreateRequest(
		sessionUUID,
		"1.1.1.1",
		"PUT",
		"http://example.com/foo1",
		[]byte("foobar"),
		map[string]string{"bbb": "foo", "aaa": "bar"},
	)

	<-time.After(time.Millisecond * 5)

	request2UUID, _ := s.CreateRequest(
		sessionUUID,
		"2.2.2.2",
		"PUT",
		"http://example.com/foo2",
		[]byte("foobar"),
		nil,
	)

	<-time.After(time.Millisecond * 5)

	request3UUID, _ := s.CreateRequest(
		sessionUUID,
		"3.3.3.3",
		"PUT",
		"http://example.com/foo3",
		[]byte("foobar"),
		map[string]string{"aaa": "bar"},
	)
	request1, _ := s.GetRequest(sessionUUID, request1UUID)
	request2, _ := s.GetRequest(sessionUUID, request2UUID)
	request3, _ := s.GetRequest(sessionUUID, request3UUID)

	req = mux.SetURLVars(req, map[string]string{"sessionUUID": sessionUUID})

	handler.ServeHTTP(rr, req)

	runtime.Gosched()
	<-time.After(time.Millisecond) // goroutine must be done

	contentBase64 := base64.StdEncoding.EncodeToString([]byte("foobar"))

	assert.JSONEq(t, `[{
		"client_address":"1.1.1.1",
		"content_base64":"`+contentBase64+`",
		"created_at_unix":`+strconv.FormatInt(request1.CreatedAt().Unix(), 10)+`,
		"headers":[{"name": "aaa", "value": "bar"},{"name": "bbb", "value": "foo"}],
		"method":"PUT",
		"url":"http://example.com/foo1",
		"uuid":"`+request1.UUID()+`"
	},{
		"client_address":"2.2.2.2",
		"content_base64":"`+contentBase64+`",
		"created_at_unix":`+strconv.FormatInt(request2.CreatedAt().Unix(), 10)+`,
		"headers":[],
		"method":"PUT",
		"url":"http://example.com/foo2",
		"uuid":"`+request2.UUID()+`"
	},{
		"client_address":"3.3.3.3",
		"content_base64":"`+contentBase64+`",
		"created_at_unix":`+strconv.FormatInt(request3.CreatedAt().Unix(), 10)+`,
		"headers":[{"name": "aaa", "value": "bar"}],
		"method":"PUT",
		"url":"http://example.com/foo3",
		"uuid":"`+request3.UUID()+`"
	}]`, rr.Body.String())
}
