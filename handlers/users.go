package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

func (h *Handlers) Users(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("users", "users")

	err := h.App.Render.JetPage(w, r, "users", vars, nil)
	if err != nil {
		h.App.ErrorLog.Error("Error rendering : ", err)
	}
}
