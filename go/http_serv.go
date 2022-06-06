package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Connected_Status struct {
	User           string
	User_Hased     string
	Rank_Id        string
	Rank_Id_Hashed string
}

var Connected Connected_Status
var templatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return to [Select_Page] ↓

//Return to page Selected need Path (string)
func Return_To_Page(w http.ResponseWriter, r *http.Request, Path string) {
	template.Must(template.ParseFiles(filepath.Join(templatesDir, Path))).Execute(w, " ")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	Return_To_Page(w, r, "../static/templates/managed_pages/404.html")
	// http.Error(w, "Method is not supported.", http.StatusBadRequest) //<-- Print [error] Method is not supported
	// fmt.Fprint(w, http.StatusBadRequest)

}

//#------------------------------------------------------------------------------------------------------------# ↓ Logout ↓

//Logout
func logout(w http.ResponseWriter, r *http.Request) {
	del(w, r)
	fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Login ↓

//Login Page
func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		Return_To_Page(w, r, "../static/templates/login.html")

	} else if r.Method == "POST" {

		r.ParseForm()

		var (
			mail_login = r.Form["mail_login"][0]
			pswd_login = r.Form["password_login"][0]
			Check      = Check_If_Exist(mail_login, pswd_login, "Pswd", "user", "Login")
		)

		if Check == true {
			SettCookie(w, r) //send cookie first
			fmt.Fprint(w, `<script> window.alert('Your are connected') </script>`)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
			//--> redirect to index

			//--> redirect to index.html
		} else { // <-- Send Error

			fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
			Return_To_Page(w, r, "../static/templates/login.html")

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
		Return_To_Page(w, r, "../static/templates/register.html")

	} else if r.Method == "POST" {

		r.ParseForm()

		var (
			User_Register         = r.Form["username_register"][0]
			Email_Register        = r.Form["email_register"][0]
			Pswd_Register         = r.Form["password_register"][0]
			Pswd_Register_Confirm = r.Form["password_register_confirm"][0]
			Check_User            = Check_If_Exist(User_Register, "", "Name", "user", "Register")
			Check_Email           = Check_If_Exist(Email_Register, "", "Email", "user", "Register")
			Hash_Pswd             = initHashPswd(Pswd_Register)
		)

		if Check_User == true && Check_Email == true && Pswd_Register == Pswd_Register_Confirm { // <-- If all is ok

			ADD_User_To_BDD(User_Register, Hash_Pswd, Email_Register, "3") // <-- Add to bdd & hash pswd
			fmt.Fprint(w, `<script> window.alert('Login now') </script>`)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)

			//--> redirect to index
			return

		} else { // <-- Check the wrong selection

			error_message := ""
			if Pswd_Register == Pswd_Register_Confirm {
				if Check_Email == false && Check_User == false {
					error_message = "email and username are"
				} else if Check_Email == false {
					error_message = "email is"
				} else if Check_User == false {
					error_message = "username is"
				}
				fmt.Fprint(w, "<script> window.alert('This "+error_message+" already in use, try again'); </script>")
			} else {
				fmt.Fprint(w, "<script> window.alert('Bad password confirmation, try again'); </script>")
			}

			Return_To_Page(w, r, "../static/templates/register.html")
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
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/admin_panel", Admin_Panel)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
