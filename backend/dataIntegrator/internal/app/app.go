package app

import "fmt"

type Deps struct {
}

type App struct {
}

func New(d *Deps) *App { //nolint: revive
	return &App{}
}

func (a *App) Start() {
	fmt.Println("F*ck The Industry!!") //nolint: forbidigo
}
