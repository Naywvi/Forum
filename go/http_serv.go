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
		user_exist := CheckIfExistLogin(r.Form["mail_login"][0], r.Form["password_login"][0], initHashPswd(r.Form["password_login"][0])) //<-- check witch hash pswd
		if user_exist == true {
			fmt.Fprint(w, "Connection > redirection + token")
		} else {
			fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/login.html"))).Execute(w, " ")
		}
	} else {
		checkHttpError(w, r)
		return
	}

}

//#------------------------------------------------------------------------------------------------------------# ↓ Register ↓ //+check if not exist in bdd
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../templates/register.html"))).Execute(w, " ")
	} else if r.Method == "POST" {
		r.ParseForm()
		user_not_exist := CheckIfExist(r.Form["username_register"][0], "name", "user")
		mail_not_exist := CheckIfExist(r.Form["email_register"][0], "Email", "user")
		if user_not_exist == true && mail_not_exist == true {
			fmt.Fprint(w, "Enregistrer > redirection + token")
			ADDUserToBDD(r.Form["username_register"][0], initHashPswd(r.Form["password_register"][0]), r.Form["email_register"][0]) //<-- Add hash pswd
		} else {
			if mail_not_exist == false {
				fmt.Fprint(w, "<script> window.alert('This email is already in use, try again'); </script>")
			} else if user_not_exist == false {
				fmt.Fprint(w, "<script> window.alert('This username is already in use, try again'); </script>")
			}
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
