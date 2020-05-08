package rest

import (
	"io"
	"net/http"
	"strings"

	"github.com/gidor/ube/pkg/apicfg"
)

// NewHealthCheckHandler add route for healthcheck
func (h *Handler) AddApi() {

	// h.router.HandleFunc("/pippo", pippo).Methods("GET")
	cfg, err := apicfg.GetApiCfg()
	if err != nil {
		h.logger.Log.Error("errorereading cfg")
		h.logger.Log.Error(err.Error())
		panic(1)
	}
	print(cfg)
	for i := 0; i < len(cfg.Api); i++ {
		method := cfg.Api[i]
		var verbs []string
		if len(method.Verbs) == 0 {
			verbs = append(verbs, "GET")
		}
		for _, verb := range method.Verbs {
			verb = verbnormalize(verb)
			if verb != "" {
				verbs = append(verbs, verb)
			}
		}
		// h.router.Handle("/"+cfg.Name+"/"+method.Path, apiHandler(h.) )
		h.router.HandleFunc("/"+cfg.Name+"/"+method.Path, method.Handler).Methods(verbs...)
		h.router.HandleFunc("/"+cfg.Name+"/"+method.Path, optionHandlerFunc(verbs)).Methods("OPTIONS")
	}

}

func verbnormalize(verb string) string {
	// ensure verbs are normalized
	switch strings.ToUpper(verb) {
	case "DELETE":
		return "DELETE"
	case "GET":
		return "GET"
	// case "OPTIONS":
	// 	return "OPTIONS"
	case "PATCH":
		return "PATCH"
	case "POST":
		return "POST"
	default:
		return ""
	}
}

func in(a []string, c string) bool {
	// a simple search function
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}

func optionHandlerFunc(verbs []string) http.HandlerFunc {
	stverbs := strings.Join(verbs, ", ")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//handle preflight in here
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Authorization")
			w.Header().Set("Allow", stverbs)
			w.Header().Set("Access-Control-Allow-Methods", stverbs)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func apiHandler(h http.Handler, method apicfg.ApiMethod) http.HandlerFunc {
	verbs := []string{"OPTIONS"}
	if len(method.Verbs) == 1 {
		verbs = append(verbs, "GET")
	}
	for _, verb := range method.Verbs {
		verb = verbnormalize(verb)
		if verb != "" {
			verbs = append(verbs, verb)
		}
	}
	stverbs := strings.Join(verbs, ", ")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//handle preflight in here
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Authorization")
			w.Header().Set("Allow", stverbs)
			w.Header().Set("Access-Control-Allow-Methods", stverbs)
			w.WriteHeader(http.StatusOK)
		} else if in(verbs, r.Method) {
			method.Handler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		h.ServeHTTP(w, r)
	}
}

func pippo(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
