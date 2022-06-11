package main

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func Create_Post(w http.ResponseWriter, r *http.Request) {
	type Statement_of_user struct {
		User string
		Rank string
	}
	//<<< --- Check rank

	var (
		_, statement, User = Check_Cookie(w, r)
		pos                = Statement_of_user{}
	)
	pos.User = User
	pos.Rank = statement
	if r.Method == "GET" {
		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/create_post.html"))).Execute(w, pos)
		Return_To_Page(w, r, "../static/templates/create_post.html")
	} else if r.Method == "POST" {
		r.ParseForm()
		var (
		//Titre := r.Form["Post_Title"][0]
		//Categorie := r.FormValue("Category_Post")
		//fmt.Println(Categorie)
		//fmt.Println(Titre)

		)

	} else {

		Send_Error(w, r)

		return
	}
}
