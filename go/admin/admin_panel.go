package admin

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	Config "forum/config"
	Database "forum/database"
	User "forum/user"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Page Admin_Panel ↓

func Admin_Panel(w http.ResponseWriter, r *http.Request) {
	var (
		I_I       = Config.Instance_of_instance{}
		db, err   = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		query     = r.FormValue("") // <-- query recuperation
		all_table = []string{"categorie", "email_owner", "post", "user"}
	)
	_, _, name := User.Check_Cookie(w, r)
	I_I.Name = name
	if r.Method == "GET" { // cookie verification to go to admin panel (security)

		s, rank, _ := User.Check_Cookie(w, r)
		if s && rank == "1" {
			Return_With_Value_Admin(w, r, I_I)
		} else { //<-- rank == 4
			Config.Send_Error(w, r)
		}

	} else if r.Method == "POST" { // We accept only methods "POST" and "GET" for security
		_, _, I_I.Name = User.Check_Cookie(w, r)
		if err != nil {
			log.Fatal(err)
		}

		if query == "Create_user" {

			r.ParseForm() // <-- obligatory for receiving values from our "query"

			var ( // <-- user_input recuperation in admin_panel.html for creating new user bby this page
				username       = r.Form["username_create_by_admin"][0]
				email          = r.Form["Email_create_by_admin"][0]
				rank           = r.Form["Rank_id_create_by_admin"][0]
				pswd           = r.Form["password_create_by_admin"][0]
				Check_username = User.Check_If_Exist(username, "", "Name", "user", "Register")
				Check_email    = User.Check_If_Exist(email, "", "Email", "user", "Register")
				pswd_hash      = Database.InitHashPswd(pswd)
			)

			if Check_username && Check_email { // <-- verify if the user we want to create do not exist

				User.ADD_User_To_BDD(username, pswd_hash, email, rank)
				_, _, I_I.Name = User.Check_Cookie(w, r)
				Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
				fmt.Fprint(w, "<script> window.alert('"+username+" successfully CREATED !'); </script>")

			} else { // <-- if wrong selection

				error_message := ""

				if !Check_email && !Check_username {
					error_message = "email and username are"
				} else if !Check_email {
					error_message = "email is"
				} else if !Check_username {
					error_message = "username is"
				}

				fmt.Fprint(w, "<script> window.alert('This "+error_message+" already in use, try again'); </script>")
				_, _, I_I.Name = User.Check_Cookie(w, r)
				Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
			}

		} else if query == "See_Table" {

			I_I = Database.Select_All_Rows_Table(db, all_table)
			_, _, I_I.Name = User.Check_Cookie(w, r)
			Return_With_Value_Admin(w, r, I_I)
			return

		} else if query == "Backup" {

			I_I = Database.Select_All_Rows_Table(db, all_table)
			Set_Backup(I_I)
			_, _, I_I.Name = User.Check_Cookie(w, r)
			fmt.Fprint(w, "<script>alert('Backup create Succesfully')</script>")
			Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")

		} else if query == "Manage_Users" {

			I_I = Database.Select_All_Rows_Table(db, all_table)
			_, _, I_I.Name = User.Check_Cookie(w, r)
			Return_With_Value_Admin(w, r, I_I)
			return

		} else if query == "Create_Category" {

			var (
				Create_categorie = r.Form["create_category"][0]
				Check_categorie  = User.Check_If_Exist(Create_categorie, "", "Name", "categorie", "New_categorie")
			)

			if Check_categorie {
				Database.Inser_In_To_DB(db, "'"+Create_categorie+"'", "categorie", "Name")
			} else {
				fmt.Fprint(w, "<script> window.alert('This categorie exist.'); </script>")
			}

		} else if query == "Alert" {
			var (
				table = []string{"user"}
				it    = Config.Instance_Bdd{}
				to    = []string{}
			)
			instance := Database.Select_All_Rows_Table(db, table)
			for _, i := range instance.I {

				it.I = append(it.I, i.I...)
			}
			for _, k := range it.I {
				to = append(to, k.User.Email)
			}
			path := "../static/templates/smtp/alert.html"
			User.Alert_Smtp(to, path) //<-- Send mail to all mails of users
			fmt.Fprint(w, "<script> window.alert('Alert sent.'); </script>")
		} else if query[:14] == "Delete_account" {
			User.Delete_Account(query[14:])
			Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
			fmt.Fprint(w, "<script> window.alert('"+query[14:]+"'s accound successfully DELETED !'); </script>")
		}

	} else {
		Config.Send_Error(w, r)
		return
	}

}
func Return_With_Value_Admin(w http.ResponseWriter, r *http.Request, I Config.Instance_of_instance) {

	template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/admin/panel_admin.html"))).Execute(w, I)
}
