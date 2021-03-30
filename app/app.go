package app

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/shortdaddy0711/goTodo/model"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var store = sessions.NewCookieStore([]byte(goDotEnvVariable("SESSION_KEY")))
var rd *render.Render = render.New()

func goDotEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return os.Getenv(key)
}

func getSessionID(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	if val == nil {
		return ""
	}
	return val.(string)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	w.Header().Set("Content-Type", "application/json")
	list := model.GetTodos(sessionId)
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
	sessionId := getSessionID(r)
	w.Header().Set("Content-Type", "application/json")
	todo := model.AddTodo(r, sessionId)
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

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if strings.Contains(r.URL.Path, "/signin") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

func MakeHandler() http.Handler {

	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(CheckSignin),
		negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/api/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/api/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/api/todos/{id}", getTodoHandler).Methods("GET")
	r.HandleFunc("/api/todos/{id}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/api/complete/{id}", completeTodoHandler).Methods("GET")
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleOAuthCallBack)

	return n
}
