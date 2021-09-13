package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error
var templates = template.Must(template.ParseFiles("home.html"))

func signupPage(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	confirmpassword := req.FormValue("confirmPassword")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	birthdate := req.FormValue("birthdate")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	if password == confirmpassword {
		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.1", 500)
				return
			}

			_, err = db.Exec("INSERT INTO users(username, password, firstname, lastname,birthdate) VALUES(?, ?, ?, ?,?)", username, hashedPassword, firstname, lastname, birthdate)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.2", 500)
				return
			}

			http.Redirect(res, req, "home", 301)
			return

		case err != nil:
			http.Error(res, "Server error, unable to create your account.3", 500)
			return
		default:
			http.Redirect(res, req, "/", 301)
		}
	} else {
		res.Write([]byte("Passwords do not match"))
	}
}

func signinPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signin.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/signin", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/signin", 301)
		return
	}

	http.Redirect(res, req, "home", 301)
}

func homePage(res http.ResponseWriter, req *http.Request) {

	//Getting data from the database
	rows, err := db.Query("SELECT username, firstname, lastname, birthdate FROM users")

	if err != nil {
		http.Error(res, "Server error, unable to get data from the database1", 500)
		return
	}

	user := User{}
	users := []User{}

	//Filling a arr with the users
	for rows.Next() {
		var username, firstname, lastname, birthdate string
		err = rows.Scan(&username, &firstname, &lastname, &birthdate)
		

		if err != nil {
			http.Error(res, "Server error, unable to get data from the database2", 500)
			return
		}
		
		user.Username = username
		user.Firstname = firstname
		user.Lastname = lastname
		user.Birthdate = birthdate
		users = append(users, user)
	}

	//executing the html with a fill format
	if err := templates.Execute(res, users); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

type User struct {
	Username, Lastname, Firstname, Birthdate string
}

func main() {
	db, err = sql.Open("mysql", "root:@/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signin", signinPage)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/", signinPage)
	http.ListenAndServe(":8080", nil)
}
