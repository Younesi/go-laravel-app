package main

import (
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"

	"github.com/younesi/atlas"
)

type application struct {
	App        *atlas.Atlas
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	atlas := initApplication()
	atlas.App.ListenAndServe()
}
