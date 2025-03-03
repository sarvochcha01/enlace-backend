package main

import "github.com/sarvochcha01/enlace-backend/cmd/app"

func main() {
	app := app.App{}
	app.Initialise()
	app.Run()
}
