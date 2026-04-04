package main

import (
	"github.com/rachmanzz/fiber-starter/bootstrap"
)

func main() {
	app := bootstrap.NewApplication()
	app.Bootstrap()
	app.Run()
}
