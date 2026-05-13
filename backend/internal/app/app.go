package app

import "fmt"

type Deps struct {
}

type App struct {
}

func New(d *Deps) *App {
	return &App{}
}

func (a *App) Start() {
	fmt.Println("F*ck The Industry!!")
}
