package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"


var db *sql.DB
var err error

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

	err := db.QueryRow("SELECT username FROM userswithdate WHERE username=?", username).Scan(&user)

	if password == confirmpassword{
		switch{
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.1", 500)
				return
			}

			_, err = db.Exec("INSERT INTO userswithdate(username, password, firstname, lastname,birthdate) VALUES(?, ?, ?, ?,?)", username, hashedPassword,firstname,lastname,birthdate)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.2", 500)
				return
			}

			res.Write([]byte("User created!"))
		return
		
		case err != nil:
			http.Error(res, "Server error, unable to create your account.3", 500)
		return
		default:
			http.Redirect(res, req, "/", 301)
		}
	}else{
		res.Write([]byte("Passwords do not match"))
	}
}


func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
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
	//http.HandleFunc("/signin", loginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}