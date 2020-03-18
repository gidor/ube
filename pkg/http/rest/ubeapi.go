package rest

import (
	"io"
	"net/http"

	"github.com/gidor/ube/pkg/apicfg"
)

// NewHealthCheckHandler add route for healthcheck
func (h *Handler) AddApi() {

	// h.router.HandleFunc("/pippo", pippo).Methods("GET")
	cfg, err := apicfg.GetApiCfg()
	if err != nil {
		h.logger.Log.Error("errorereading cfg")
	}
	print(cfg)

}

func pippo(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
