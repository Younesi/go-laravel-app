package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/younesi/atlas"
	"myapp/data"
	"net/http"
)

type Handlers struct {
	App    *atlas.Atlas
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "message"
	h.App.Session.Put(r.Context(), myData, "Hello from a session!")
	msg := h.App.Session.GetString(r.Context(), myData)

	vars := make(jet.VarMap)
	vars.Set(myData, msg)

	err := h.App.Render.JetPage(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error rendering : ", err)
	}
}
