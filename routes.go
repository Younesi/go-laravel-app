package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middlewares must come before any routes

	// Auth
	a.App.Routes.Get("/auth/login", a.Handlers.Login)
	a.App.Routes.Post("/auth/login", a.Handlers.PostLogin)
	a.App.Routes.Get("/auth/logout", a.Handlers.Logout)

	// Home
	a.App.Routes.Get("/", a.Handlers.Home)
	a.App.Routes.Get("/go-page", a.Handlers.GoPage)
	a.App.Routes.Get("/jet-page", a.Handlers.JetPage)
	a.App.Routes.Get("/sessions", a.Handlers.SessionTest)

	// Users
	a.App.Routes.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println("Error in fetching users")
			return
		}

		for _, user := range users {
			a.App.InfoLog.Println(user)
			fmt.Fprintf(w, "%s : %d", user.FirstName, user.ID)
		}
	})

	//a.App.Routes.Get("/test-user-create", func(w http.ResponseWriter, r *http.Request) {
	//	u := data.User{
	//		FirstName: "Mahdi",
	//		LastName:  "Younesi",
	//		Email:     "mehdi.younesi7@gmail.com",
	//		IsActive:  1,
	//		Password:  "simple",
	//	}
	//
	//	id, err := a.Models.Users.Insert(u)
	//	if err != nil {
	//		a.App.ErrorLog.Println(err)
	//		return
	//	}
	//
	//	fmt.Fprintf(w, "%s : %d", u.FirstName, id)
	//})

	// add static files
	filServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", filServer))

	return a.App.Routes
}
