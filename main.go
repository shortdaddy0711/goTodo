package main

import (
	"net/http"

	"github.com/shortdaddy0711/goTodo/app"
	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	http.ListenAndServe(":3000", n)
}
