package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Show post ↓
func Show_Post(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Post struct {
		Id                  string
		Title               string
		Content             string
		Posted_user         string
		Posted_user_rank    string
		Likes               string
		Last_Posted         string
		Nb_Reply            string
		User_connected      string
		User_connected_rank string
	}
	type Statement_of_user struct {
		User    string
		Rank    string
		Post_Id string
		Post    []Post
	}

	var (
		_, statement, User = Check_Cookie(w, r)
		pos                = Statement_of_user{}
		rows               = Select_column("post", "Id", query)
		instance           all_bd
		POST               = Post{}
	)
	pos.User = User
	pos.Rank = statement
	pos.Post_Id = query
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

		//flemme intense
		result := ""
		flemme := Select_column("user", "name", POST.Posted_user)
		for flemme.Next() {
			errf := flemme.Scan(&instance.User.Id, &instance.User.Name, &instance.User.Pswd, &instance.User.Desc, &instance.User.Email, &instance.User.Profile_Picture, &instance.User.Rank_id)
			if errf != nil {
				log.Fatal(errf)
			}
			result = strconv.Itoa(instance.User.Rank_id)

		}
		POST.Posted_user_rank = result
		//flemme intense

		POST.Likes = instance.Post.Likes
		POST.Last_Posted = instance.Post.Last_Posted
		POST.Nb_Reply = instance.Post.Nb_Reply
		POST.User_connected = User
		POST.User_connected_rank = statement

		//<<Append the post
		pos.Post = append(pos.Post, POST)

	}

	if r.Method == "GET" {

		//<<< add post
		template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/post.html"))).Execute(w, pos)

	} else if r.Method == "POST" {
		// r.ParseForm()
		// if query == "send" {
		// 	var (
		// 		db, err  = sql.Open(Bdd.Langage, Bdd.Name)
		// 		Title    = r.Form["Post_Title"][0]
		// 		Content  = r.Form["Post_Content"][0]
		// 		cat      = r.Form["categorie_id"][0]
		// 		time     = time.Now()
		// 		timestr  = time.String()
		// 		var_p    = []string{"'" + cat + "','" + Title + "','" + Content + "','0','" + User + "','" + timestr[0:10] + "','0'"}
		// 		var_pstr = strings.Join(var_p, "")
		// 	)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	Inser_In_To_DB(db, var_pstr, "post", Extract_File("../bdd/post_table.sql", 11, 12)) //<-- Redirect to post
		// 	fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
	} else {
		Send_Error(w, r)

		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Create post ↓

//Creat post func
func Create_Post(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
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
	if statement != "4" {

		if r.Method == "GET" {
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/create_post.html"))).Execute(w, pos)

		} else if r.Method == "POST" {
			r.ParseForm()
			if query == "send" {
				var (
					db, err  = sql.Open(Bdd.Langage, Bdd.Name)
					Title    = r.Form["Post_Title"][0]
					Content  = r.Form["Post_Content"][0]
					cat      = r.Form["categorie_id"][0]
					time     = time.Now()
					timestr  = time.String()
					var_p    = []string{"'" + cat + "','" + Title + "','" + Content + "','0','" + User + "','" + timestr[0:10] + "','0'"}
					var_pstr = strings.Join(var_p, "")
				)
				if err != nil {
					log.Fatal(err)
				}
				Inser_In_To_DB(db, var_pstr, "post", Extract_File("../bdd/post_table.sql", 11, 12)) //<-- Redirect to post
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			} else {
				Send_Error(w, r)

				return
			}

		} else {

			Send_Error(w, r)

			return
		}

	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
	}
}
