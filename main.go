package main

import (
	"log"

	"SweetDreams/app"
	"SweetDreams/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := config.NewConfig()
	app.ConfigAndRunApp(config)
}
