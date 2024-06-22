package main

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// Middleware

	a.App.Routes.Use(a.Middleware.CheckRemember)

	// Auth
	a.App.Routes.Get("/auth/login", a.Handlers.Login)
	a.App.Routes.Post("/auth/login", a.Handlers.PostLogin)
	a.App.Routes.Get("/auth/logout", a.Handlers.Logout)
	a.App.Routes.Get("/auth/forgot-password", a.Handlers.Forgot)
	a.App.Routes.Post("/auth/forgot-password", a.Handlers.PostForgot)

	// Home
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/go-page", a.Handlers.GoPage)
	a.App.Routes.Get("/jet-page", a.Handlers.JetPage)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	// Forms
	a.App.Routes.Get("/form", a.Handlers.Form)
	a.App.Routes.Post("/form", a.Handlers.SubmitForm)

	a.App.Routes.Get("/json-test", a.Handlers.JsonTest)
	a.App.Routes.Get("/download-test", a.Handlers.DownloadFileTest)
	a.App.Routes.Get("/encryption-test", a.Handlers.CryptoTest)
	a.App.Routes.Get("/cache-test", a.Handlers.CacheTest)

	// Users
	a.App.Routes.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Error("error fetching users")
			return
		}

		for _, user := range users {
			fmt.Fprintf(w, "%s : %d", user.FirstName, user.ID)
		}
	})

	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "Mahdi",
			LastName:  "Younesi",
			Email:     "mehdi.younesi7@gmail.com",
			IsActive:  1,
			Password:  "simple",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Error("error inserting user", err)
			return
		}

		fmt.Fprintf(w, "%s : %d", u.FirstName, id)
	})

	a.App.Routes.Get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Error("error fetching users", err)
			return
		}
		for _, x := range users {
			fmt.Fprint(w, x.LastName)
		}
	})

	a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Error("error fetching user", err)
			return
		}

		fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	})

	a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Error("error fetching user", err)
			return
		}

		u.LastName = a.App.RandomString(10)
		err = u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Error("error updating user", err)
			return
		}

		fmt.Fprintf(w, "updated last name to %s", u.LastName)

	})

	// add static files
	filServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", filServer))

	return a.App.Routes
}
