package api

import (
	"net/http"
	mw "profile/middleware"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/profiles", h.HandleProfileCreation)

	mux.HandleFunc("GET /api/profiles/{uuid}", h.HandleProfileRetrievalByID)

	mux.HandleFunc("GET /api/profiles", h.HandleAllProfileRetrievalWithFilter)

	mux.Handle("DELETE /api/profiles/{uuid}", mw.BearerAuth(http.HandlerFunc(h.HandleProfileDeletionByID)))

	mux.Handle("POST /admin", mw.BearerAuth(http.HandlerFunc(h.HandleAdmin)))

	return mux
}
