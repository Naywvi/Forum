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

//#------------------------------------------------------------------------------------------------------------# ↓ Home Page ↓

//Home page
func forum(w http.ResponseWriter, r *http.Request) {
	type Statement_of_user struct {
		User string
		Rank string
	}
	if r.Method == "GET" {
		//<<< --- Check rank

		var (
			_, statement, User = Check_Cookie(w, r)
			pos                = Statement_of_user{}
		)
		pos.User = User
		pos.Rank = statement

		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/forum.html"))).Execute(w, pos)

		//<<< --- Check rank
		Return_To_Page(w, r, "../static/templates/forum.html")

	} else if r.Method == "POST" {

	} else {

		Send_Error(w, r)

		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓

//Server Http
func httpServ() {
	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)

	//<<< register_login_logout.go
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/resend_mail", Resend_Mail)
	http.HandleFunc("/valide_password", valide_password_page)
	http.HandleFunc("/validation_mail", Validation_URLbyMail)
	http.HandleFunc("/reset_password_page", reset_password_page)
	//<<<

	//<<< profil.go
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/profil/edit", edit_desc)
	http.HandleFunc("/delete-account", profildeleted)
	//<<<

	//<<< post.go
	http.HandleFunc("/post", Show_Post)
	http.HandleFunc("/create_post", Create_Post)
	//<<<

	//<<< admin_panel.go
	http.HandleFunc("/admin_panel", Admin_Panel)
	//<<<

	//<<< categorie.go
	http.HandleFunc("/categorie", Show_Categorie)
	//<<<
	http.HandleFunc("/forum", forum)

	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
