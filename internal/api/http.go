package api

import "net/http"

func (h *Handler) Account(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("您好!!!"))
	w.WriteHeader(200)
}
