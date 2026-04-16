package api

import (
	"net/http"
)

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/profiles", h.HandleProfileCreation)
	
	mux.HandleFunc("GET /api/profiles/{id}", h.HandleProfileRetrievalByID)

	return mux
}
