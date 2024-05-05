package main

import (
	"log"
	"myapp/data"
	"myapp/handlers"
	"os"

	"github.com/younesi/atlas"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init atlas
	cel, err := atlas.New(path)

	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "MyApp"
	cel.InfoLog.Println("Debug: ", cel.Debug)

	myHandlers := &handlers.Handlers{
		App: cel,
	}

	app := &application{
		App:      cel,
		Handlers: myHandlers,
	}
	app.App.Routes = app.routes()
	app.Models = data.New(app.App.DB.Pool) // sql db
	myHandlers.Models = app.Models

	return app
}
