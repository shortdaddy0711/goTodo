package main

import (
	"log"
	"net/http"

	"github.com/shortdaddy0711/goTodo/app"
)

func main() {
	m := app.MakeHandler()

	log.Println("started app")
	err := http.ListenAndServe(":8000", m)
	if err != nil {
		log.Fatal(err)
	}
}
