package handlers

import (
	"net/http"
	"sync/atomic"
)

// readyz is a readiness probe.
func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("isReady: " + isReady)
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
