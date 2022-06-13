package categorie

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"

	Config "forum/config"
	Database "forum/database"
	User "forum/user"
)

func Show_Categorie(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Post struct {
		Id          string
		Title       string
		Content     string
		Posted_user string
		Likes       string
		Last_Posted string
		Nb_Reply    string
	}

	type Statement_of_user struct {
		Username  string
		Rank      string
		Categorie string
		Post      []Post
	}
	var (
		_, statement, Username = User.Check_Cookie(w, r)
		pos                    = Statement_of_user{}
		POST                   = Post{}
		rows                   = Database.Select_column("post", "Id_cat", query)
		instance               Config.All_bd
		Check_categorie        = User.Check_If_Exist(query, "", "Name", "categorie", "New_categorie")
	)

	pos.Username = Username
	pos.Rank = statement
	pos.Categorie = query

	//<<< add post
	for rows.Next() {
		err := rows.Scan(&instance.Post.Id, &instance.Post.Id_cat, &instance.Post.Title_post, &instance.Post.Content, &instance.Post.Likes, &instance.Post.Posted_user, &instance.Post.Last_Posted, &instance.Post.Nb_Reply)

		if err != nil {
			log.Fatal(err)
		}
		POST.Id = strconv.Itoa(instance.Post.Id)
		POST.Title = instance.Post.Title_post
		POST.Content = instance.Post.Content
		POST.Posted_user = instance.Post.Posted_user
		POST.Likes = instance.Post.Likes
		POST.Last_Posted = instance.Post.Last_Posted
		POST.Nb_Reply = instance.Post.Nb_Reply
		pos.Post = append(pos.Post, POST)

	}
	//<<< add post
	if !Check_categorie {

		if r.Method == "GET" { //Besoin d'une instance de tout mes posts de cette catégorie
			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/list_post.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

		} else {

			Config.Send_Error(w, r)

			return
		}
	} else {
		Config.Send_Error(w, r)

		return
	}
}
