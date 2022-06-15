package user

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	Config "forum/config"
)

func Login(w http.ResponseWriter, r *http.Request) { //Login Page
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

	//<<< --- Check rank
	if pos.Rank == "4" {
		if r.Method == "GET" {

			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/login.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

			r.ParseForm()

			var (
				mail_login = r.Form["mail_login"][0]
				pswd_login = r.Form["password_login"][0]
				Check      = Check_If_Exist(mail_login, pswd_login, "Pswd", "user", "Login")
			)
			if Check {
				SettCookie(w, r) //send cookie first
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
				//--> redirect to index

				//--> redirect to index.html
			} else { // <-- Send Error

				fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
				template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/login.html"))).Execute(w, pos)

			}

		} else { // <-- If r.Method != Get/Post

			Config.Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}
