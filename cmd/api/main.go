package main

import (
	"log"

	"github.com/radityarps/gophersocial/internal/env"
)

func main() {
	app := &application{
		config: config{
			addr: env.GetString("ADDR", ":8080"),
		},
	}
	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
