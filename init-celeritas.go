package main

import (
	"log"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"

	"github.com/younesi/atlas"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init atlas
	at, err := atlas.New(path)

	if err != nil {
		log.Fatal(err)
	}

	at.AppName = "MyApp"
	at.InfoLog.Info("Debug mode:", "debug", at.Debug)

	myMiddleware := &middleware.Middleware{
		App: at,
	}
	myHandlers := &handlers.Handlers{
		App:        at,
		Middleware: myMiddleware,
	}
	myModels := data.New(at.DB.Pool)

	myHandlers.Models = myModels
	myMiddleware.Models = myModels

	app := &application{
		App:        at,
		Handlers:   myHandlers,
		Models:     myModels,
		Middleware: myMiddleware,
	}
	app.App.Routes = app.routes()

	return app
}
