package main

import (
	"github.com/jst-Frenzy/iot/backend/internal/app"
)

func main() {
	application := app.New(&app.Deps{})

	application.Start()
}
