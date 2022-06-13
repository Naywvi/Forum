package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	Admin "forum/admin"
	Categorie "forum/categories"
	Config "forum/config"
	Post "forum/post"
	User "forum/user"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓

//Server Http
func HttpServ() {
	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)

	//<<< register_login_logout.go
	http.HandleFunc("/login", User.Login)
	http.HandleFunc("/register", User.Register)
	http.HandleFunc("/logout", User.Logout)
	http.HandleFunc("/resend_mail", User.Resend_Mail)
	http.HandleFunc("/valide_password", User.Valide_password_page)
	http.HandleFunc("/validation_mail", User.Validation_URLbyMail)
	http.HandleFunc("/reset_password_page", User.Reset_password_page)
	//<<<

	//<<< profil.go
	http.HandleFunc("/profil", User.Profil)
	http.HandleFunc("/profil/edit", User.Edit_desc)
	http.HandleFunc("/delete-account", User.Profildeleted)
	//<<<

	//<<< post.go
	http.HandleFunc("/post", Post.Show_Post)
	http.HandleFunc("/create_post", Post.Create_Post)
	//<<<

	//<<< admin_panel.go
	http.HandleFunc("/admin_panel", Admin.Admin_Panel)
	//<<<

	//<<< categorie.go
	http.HandleFunc("/categorie", Categorie.Show_Categorie)
	//<<<
	http.HandleFunc("/forum", forum)

	fmt.Println("Started https serv successfully on http://localhost:8080")
	fmt.Print(http.ListenAndServe(":1010", nil))

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
			_, statement, User = User.Check_Cookie(w, r)
			pos                = Statement_of_user{}
		)
		pos.User = User
		pos.Rank = statement

		template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/forum.html"))).Execute(w, pos)

		//<<< --- Check rank
		Config.Return_To_Page(w, r, "../static/templates/forum.html")

	} else if r.Method == "POST" {

	} else {

		Config.Send_Error(w, r)

		return
	}
}
