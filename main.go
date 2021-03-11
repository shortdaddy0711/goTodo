package main

import (
	"log"
	"net/http"

	"github.com/shortdaddy0711/goTodo/app"
	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("started app")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
