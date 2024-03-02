// main.go
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "adv_prog"
)

var db *sql.DB

type User struct {
	Username string
	Password string
	Role     string
}

func main() {

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/register", registerHandler).Methods("GET", "POST") // Allow GET and POST requests
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/index", dashboardHandler).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Serve the registration form HTML here
		tmpl, err := template.ParseFiles("register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		ageStr := r.FormValue("age")
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password, firstname, lastname, age) VALUES ($1, $2, $3, $4, $5, $6)", username, email, password, firstname, lastname, age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if password != storedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)

}
