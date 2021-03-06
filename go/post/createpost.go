package post

import (
	"database/sql"
	"fmt"
	Config "forum/config"
	Database "forum/database"
	User "forum/user"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//Creat post func
func Create_Post(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Category struct {
		Id   int
		Name string
	}
	type Statement_of_user struct {
		User      string
		Rank      string
		Categorie []Category
	}
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
	if statement != "4" {

		if r.Method == "GET" {
			var (
				rows = Database.Select_All_From_DB(db, "categorie")
			)
			for rows.Next() {
				err := rows.Scan(&cate.Id, &cate.Name)
				if err != nil {
					log.Fatal(err)
				}

			}
			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/create_post.html"))).Execute(w, pos)

		} else if r.Method == "POST" {
			r.ParseForm()
			if query == "send" {
				var ( //Comment_content_parse_sql = strings.Replace(Comment_content_parse, "'", "`", 10000) //<< Replace ' to > `  protect from sql_exploit
					Title_html   = r.Form["Post_Title"][0]
					Title        = strings.Replace(Title_html, "'", "’", 10000) //<< Replace ' to > `  protect from sql_exploit
					Content_html = r.Form["Post_Content"][0]
					Content      = strings.Replace(Content_html, "'", "’", 10000)
					cat          = r.Form["categorie_id"][0]
					time         = time.Now()
					timestr      = time.String()
					var_p        = []string{"'" + cat + "','" + Title + "','" + Content + "','0','" + User + "','" + timestr[0:10] + "','0'"}
					var_pstr     = strings.Join(var_p, "")
				)

				Database.Inser_In_To_DB(db, var_pstr, "post", Database.Extract_File("../bdd/post_table.sql", 11, 12)) //<-- Insert new post on bdd
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			} else {
				Config.Send_Error(w, r)

				return
			}

		} else {

			Config.Send_Error(w, r)

			return
		}

	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
	}
}
