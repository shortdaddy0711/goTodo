package app

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/shortdaddy0711/goTodo/model"
	"github.com/unrolled/render"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var rd *render.Render

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list := model.GetTodos()
	rd.JSON(w, http.StatusOK, list)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	todo := model.GetTodo(id)
	rd.JSON(w, http.StatusOK, todo)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todo := model.AddTodo(r)
	rd.JSON(w, http.StatusCreated, todo)
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	ok := model.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, true)
	} else {
		rd.JSON(w, http.StatusBadRequest, false)
	}
}

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	complete := r.FormValue("complete")
	ok := model.CompleteTodo(id, complete)
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
	r.HandleFunc("/api/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/api/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/api/todos/{id}", getTodoHandler).Methods("GET")
	r.HandleFunc("/api/todos/{id}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/api/complete/{id}", completeTodoHandler).Methods("GET")

	return r
}
