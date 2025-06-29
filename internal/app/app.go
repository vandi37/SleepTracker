package app

import "context"

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) {

}
