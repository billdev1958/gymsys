package v1

import "net/http"

func (h *handler) UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/payment", h.RegisterUser)
}
