package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shortdaddy0711/goTodo/model"
	"github.com/unrolled/render"
)

var rd *render.Render

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := model.GetTodos()
	rd.JSON(w, http.StatusOK, list)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	todo := model.AddTodo(name)
	rd.JSON(w, http.StatusCreated, todo)
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intId, _ := strconv.Atoi(vars["id"])
	ok := model.RemoveTodo(intId)
	if ok {
		rd.JSON(w, http.StatusOK, true)
	} else {
		rd.JSON(w, http.StatusBadRequest, false)
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intId, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := model.CompleteTodo(intId, complete)
	if ok {
		rd.JSON(w, http.StatusOK, true)
	} else {
		rd.JSON(w, http.StatusBadRequest, false)
	}
}

func MakeHandler() http.Handler {

	rd = render.New()
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete/{id:[0-9]+}", completeTodoHandler).Methods("GET")

	return r
}
