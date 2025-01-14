package get

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/responder"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/storage"
)

func NewHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUUID, sessionFound := mux.Vars(r)["sessionUUID"]
		if !sessionFound {
			responder.JSON(w, api.NewServerError(http.StatusInternalServerError, "cannot extract session UUID"))

			return
		}

		requestUUID, requestFound := mux.Vars(r)["requestUUID"]
		if !requestFound {
			responder.JSON(w, api.NewServerError(http.StatusInternalServerError, "cannot extract request UUID"))

			return
		}

		request, gettingErr := storage.GetRequest(sessionUUID, requestUUID)
		if gettingErr != nil {
			responder.JSON(w, api.NewServerError(
				http.StatusInternalServerError, "cannot read request data: "+gettingErr.Error(),
			))

			return
		}

		if request == nil {
			responder.JSON(w, api.NewServerError(
				http.StatusNotFound, fmt.Sprintf("request with UUID %s was not found", requestUUID),
			))

			return
		}

		var (
			headersMap = request.Headers()
			headers    = make([]api.SessionRequestHeader, 0, len(headersMap))
		)

		for name, value := range headersMap {
			headers = append(headers, api.SessionRequestHeader{Name: name, Value: value})
		}

		sort.SliceStable(headers, func(j, k int) bool { return headers[j].Name < headers[k].Name })

		responder.JSON(w, api.SessionRequest{
			UUID:          request.UUID(),
			ClientAddr:    request.ClientAddr(),
			Method:        request.Method(),
			ContentBase64: base64.StdEncoding.EncodeToString(request.Content()),
			Headers:       headers,
			URI:           request.URI(),
			CreatedAtUnix: request.CreatedAt().Unix(),
		})
	}
}
