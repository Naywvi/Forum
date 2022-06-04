package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var templatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return to [Select_Page] ↓

//Return to page Selected need Path (string)
func Return_To_Page(w http.ResponseWriter, r *http.Request, Path string) {
	template.Must(template.ParseFiles(filepath.Join(templatesDir, Path))).Execute(w, " ")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method is not supported.", http.StatusBadRequest)
	fmt.Fprint(w, http.StatusBadRequest)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Login ↓

//Login Page
func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		Return_To_Page(w, r, "../templates/login.html")

	} else if r.Method == "POST" {

		r.ParseForm()

		var (
			mail_login = r.Form["mail_login"][0]
			pswd_login = r.Form["password_login"][0]
			Hash_Pswd  = initHashPswd(pswd_login)
			user_exist = Check_If_Exist_Login(mail_login, pswd_login, Hash_Pswd) //<-- check witch hash pswd in bdd
		)

		if user_exist == true {

			//
			SetCookie(w, r)
			//
			fmt.Fprint(w, "Connection > redirection + token")

		} else { // <-- Send Error

			fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
			Return_To_Page(w, r, "../templates/login.html")

		}

	} else { // <-- If r.Method != Get/Post

		Send_Error(w, r)
		return

	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Register ↓ //+check if not exist in bdd

//Register Page
func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		Return_To_Page(w, r, "../templates/register.html")

	} else if r.Method == "POST" {

		r.ParseForm()

		var (
			Check_User     = Check_If_Exist(r.Form["username_register"][0], "name", "user")
			Check_Email    = Check_If_Exist(r.Form["email_register"][0], "Email", "user")
			User_Register  = r.Form["username_register"][0]
			Email_Register = r.Form["email_register"][0]
			Pswd_Register  = r.Form["password_register"][0]
			Hash_Pswd      = initHashPswd(Pswd_Register)
		)

		if Check_User == true && Check_Email == true { // <-- If all is ok

			fmt.Fprint(w, "Enregistrer > redirection + token")
			ADD_User_To_BDD(User_Register, Hash_Pswd, Email_Register) // <-- Add to bdd & hash pswd

		} else { // <-- Check the wrong selection

			if Check_Email == false {
				fmt.Fprint(w, "<script> window.alert('This email is already in use, try again'); </script>")
			} else if Check_User == false {
				fmt.Fprint(w, "<script> window.alert('This username is already in use, try again'); </script>")
			}

			Return_To_Page(w, r, "../templates/register.html")
		}

	} else { // <-- If r.Method != Get/Post

		Send_Error(w, r)
		return

	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓

//Server Http
func httpServ() {
	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
