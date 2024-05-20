package handlers

import (
	"fmt"
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

func (h *Handlers) CryptoTest(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello, world"
	fmt.Fprint(w, "Unencrypted: "+plainText+"\n")

	encrypted, err := h.encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.ErrInternalServer(w, r)
		return
	}

	fmt.Fprint(w, "encrypted: "+encrypted+"\n")

	decrypted, err := h.decrypt(encrypted)
	if err != nil {
		h.App.ErrorLog.Println(err)
		h.App.ErrInternalServer(w, r)
		return
	}

	fmt.Fprint(w, "decrypted: "+decrypted+"\n")
}

func (h *Handlers) CacheTest(w http.ResponseWriter, r *http.Request) {
	cacheKey := "passion"
	err := h.App.Cache.Set(cacheKey, "GoLang")
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	cacheValue, err := h.App.Cache.Get(cacheKey)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	var payload struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Message string `json:"message"`
	}

	payload.Key = cacheKey
	payload.Value = cacheValue.(string)
	payload.Message = "Fetched key from cache"

	err = h.App.WriteJson(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}
