package handlers

import (
	"myapp/data"
	"myapp/middleware"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/younesi/atlas"
)

type Handlers struct {
	App        *atlas.Atlas
	Models     data.Models
	Middleware *middleware.Middleware
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

func (h *Handlers) JsonTest(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Id      int64    `json:id`
		Name    string   `json:name`
		Hobbies []string `json:hobbies`
	}

	payload.Id = 7
	payload.Name = "Mahdi"
	payload.Hobbies = []string{"Games", "Hokey", "Sweaming", "Hanging out with friends"}

	err := h.App.WriteJson(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) DownloadFileTest(w http.ResponseWriter, r *http.Request) {
	h.App.DownloadFile(w, r, "./public/images", "atlas.png")
}
