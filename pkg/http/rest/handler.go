package rest

import (
	"github.com/gidor/ube/pkg/http/rest/middleware"
	"github.com/gidor/ube/pkg/infra"
	"github.com/gorilla/mux"
)

// Handler handles http rest requests
type Handler struct {
	logger *infra.Logger
	router *mux.Router
}

// HTTPError data model for http error
type HTTPError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// CreateHandler create a new http rest handler
func CreateHandler(l *infra.Logger) *Handler {
	h := &Handler{
		logger: l,
		router: mux.NewRouter(),
	}
	mw := middleware.CreateMiddleware(h.logger)

	h.router.Use(mw.JSONResponse)
	h.router.Use(mw.Logging)

	return h
}

// GetRouter returns the router
func (h *Handler) GetRouter() *mux.Router {
	return h.router
}
