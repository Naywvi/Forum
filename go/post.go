package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

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
