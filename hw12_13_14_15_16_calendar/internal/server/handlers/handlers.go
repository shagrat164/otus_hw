package handlers

import "net/http"

func StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //nolint:revive
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
