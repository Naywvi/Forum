package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	Admin "forum/admin"
	Categorie "forum/categories"
	Config "forum/config"
	Database "forum/database"
	Post "forum/post"
	Security "forum/security"
	User "forum/user"
)

func HttpServ() { // start the http server and configure the routes
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	// redirect index.html
	// http.HandleFunc("/", Config.RedirectHome)

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
	http.HandleFunc("/forum", home)

	fmt.Println("Started https serv successfully on http://localhost:1010")

	fmt.Print(http.ListenAndServe(":1010", routehandler(http.DefaultServeMux)))

}

func routehandler(h http.Handler) http.Handler { // route handler to call the ratelimiter before routing the request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Security.Ratelimit(w, r)
		h.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) { //Home page
	type Category struct {
		Id   int
		Name string
	}
	type Statement_of_user struct {
		User      string
		Rank      string
		Categorie []Category
	}
	if r.Method == "GET" {
		//<<< --- Check rank

		var (
			db, err            = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
			_, statement, User = User.Check_Cookie(w, r)
			pos                = Statement_of_user{}
			cate               Category
			rows               = Database.Select_All_From_DB(db, "categorie")
		)

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			err := rows.Scan(&cate.Id, &cate.Name)
			if err != nil {
				log.Fatal(err)
			}
			pos.Categorie = append(pos.Categorie, cate)
		}
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
