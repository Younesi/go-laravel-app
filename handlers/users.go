package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"net/http"
)

func (h *Handlers) Users(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("users", "users")

	err := h.App.Render.JetPage(w, r, "users", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
	}
}
