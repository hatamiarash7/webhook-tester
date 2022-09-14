package delete

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/handlers/api"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/http/responder"
	"github.com/hatamiarash7/webhook-tester/internal/pkg/storage"
	jsoniter "github.com/json-iterator/go"
)

func NewHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUUID, sessionFound := mux.Vars(r)["sessionUUID"]
		if !sessionFound {
			responder.JSON(w, api.NewServerError(http.StatusInternalServerError, "cannot extract session UUID"))

			return
		}

		// delete session
		if result, err := storage.DeleteSession(sessionUUID); err != nil {
			responder.JSON(w, api.NewServerError(http.StatusInternalServerError, err.Error()))

			return
		} else if !result {
			responder.JSON(w, api.NewServerError(
				http.StatusNotFound, fmt.Sprintf("session with UUID %s was not found", sessionUUID),
			))

			return
		}

		// and recorded session requests
		if _, err := storage.DeleteRequests(sessionUUID); err != nil {
			responder.JSON(w, api.NewServerError(http.StatusInternalServerError, err.Error()))

			return
		}

		responder.JSON(w, output{Success: true})
	}
}

type output struct {
	Success bool `json:"success"`
}

func (o output) ToJSON() ([]byte, error) { return jsoniter.ConfigFastest.Marshal(o) }
