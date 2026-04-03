package main

import (
	"github.com/rachmanzz/fiber-starter/bootstrap"
	"go.uber.org/zap"
)

func main() {
	app := bootstrap.NewApplication()
	app.RegisterProviders()
	if err := app.Bootstrap(); err != nil {
		zap.L().Info("error start app")
	}
}
