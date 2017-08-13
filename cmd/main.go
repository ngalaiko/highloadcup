package main

import (
	"github.com/ngalayko/highloadcup"
	"log"
)

func main() {
	app := highloadcup.NewApp()

	if err := app.ServeHTTP(); err != nil {
		log.Panic(err)
	}
}
