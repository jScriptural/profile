package api

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/profiles", h.HandleProfileCreation)
	
	mux.HandleFunc("GET /api/profiles/{uuid}", h.HandleProfileRetrievalByID)

	mux.HandleFunc("GET /api/profiles",h.HandleAllProfileRetrievalWithFilter)

	return mux
}
