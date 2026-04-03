package bootstrap

import "github.com/rachmanzz/fiber-starter/cores"

type Application struct{}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Bootstrap() error {
	_ = cores.CreateContract()
	return nil
}
