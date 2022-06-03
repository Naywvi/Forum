package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var templatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓
func checkHttpError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method is not supported.", http.StatusBadRequest)
	fmt.Fprint(w, http.StatusBadRequest)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Login ↓
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/login.html"))).Execute(w, " ")
	} else if r.Method == "POST" {
		r.ParseForm()
		// SEND TO BDD
		fmt.Println("username:", r.Form["mail_login"])
		fmt.Println("password:", r.Form["password_login"])
	} else {
		checkHttpError(w, r)
		return
	}
	//Recup input
}

//#------------------------------------------------------------------------------------------------------------# ↓ Register ↓
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/register.html"))).Execute(w, " ")
	} else if r.Method == "POST" {
		r.ParseForm()
		user_not_exist := CheckIfExist(r.Form["username_register"][0])
		if user_not_exist == true {
			fmt.Fprint(w, "Enregistrer > redirection + token")
			SendUserToBDD(r.Form["username_register"][0], r.Form["password_register"][0], r.Form["email_register"][0])
		} else {
			fmt.Fprint(w, "<script> window.alert('This username is already in use'); </script>")
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/register.html"))).Execute(w, " ")
		}

	} else {
		checkHttpError(w, r)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓
func httpServ() {

	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
