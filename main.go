package main

import (
	"github.com/younesi/atlas"
	"myapp/data"
	"myapp/handlers"
)

type application struct {
	App      *atlas.Atlas
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
