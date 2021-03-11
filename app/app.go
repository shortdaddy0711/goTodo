package app

import (
	"net/http"
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()
	return r

}
