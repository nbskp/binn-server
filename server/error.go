package server

import (
	"net/http"

	"github.com/nbskp/binn-server/binn"
)

func handleError(w http.ResponseWriter, err error, defaultStatusCode int) {
	bErr, ok := err.(*binn.BinnError)
	if !ok {
		w.WriteHeader(defaultStatusCode)
		return
	}
	switch bErr.Code {
	case binn.CodeExpiredSubscription:
		w.WriteHeader(http.StatusBadRequest)
	case binn.CodeNotFoundSubscribedBottle:
		w.WriteHeader(http.StatusBadRequest)
	case binn.CodeExpiredBottle:
		w.WriteHeader(http.StatusBadRequest)
	case binn.CodeNotFoundBottle:
		w.WriteHeader(http.StatusBadRequest)
	case binn.CodeUnavailableBottle:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(defaultStatusCode)
	}
}
