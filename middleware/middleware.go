package middleware

import (
	"myapp/data"

	"github.com/younesi/atlas"
)

type Middleware struct {
	App    *atlas.Atlas
	Models data.Models
}
