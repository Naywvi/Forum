package user

import (
	"fmt"
	Config "forum/config"
	Database "forum/database"
	"html/template"
	"net/http"
	"path/filepath"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Reset password Page ↓
func Reset_password_page(w http.ResponseWriter, r *http.Request) {
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

			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

			r.ParseForm()

			var (
				email_to_reset = r.Form["reset_email"][0]
				email_exist    = Check_If_Exist(email_to_reset, "", "Email", "user", "Register") // true -> don't exist
			)

			if email_exist {

				template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
				fmt.Fprint(w, "<script> window.alert('Mail don't exist'); </script>")

			} else {

				var (
					email_tab       = []string{email_to_reset}
					validation_hash = Return_From_Table(email_to_reset, "user", "Reset_password")
				)

				if validation_hash == "error" || len(validation_hash) == 0 {

					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
					fmt.Fprint(w, "<script> window.alert('Contact administrator'); </script>")

				} else {

					Init_Smtp(email_tab, "", validation_hash, "Reset")
					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
					fmt.Fprint(w, "<script> window.alert('Sent ! Now check your emails.'); </script>")
				}

			}

		} else {

			Config.Send_Error(w, r)

			return
		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

func Reset_Password(New_password, Last_password string) {
	New_password = Database.InitHashPswd(New_password)
	Database.Update_Field("user", "Pswd", "Pswd", Last_password, New_password)
}
