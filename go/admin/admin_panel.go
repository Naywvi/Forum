package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	Config "forum/config"
	Database "forum/database"
	User "forum/user"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Page Admin_Panel ↓

func Admin_Panel(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" { // cookie verification to go to admin panel (security)

		s, rank, _ := User.Check_Cookie(w, r)
		if s == true && rank == "1" {
			Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
		} else { //<-- rank == 4
			Config.Send_Error(w, r)
		}

	} else if r.Method == "POST" { // We accept only methods "POST" and "GET" for security
		var (
			db, err   = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
			I_I       Config.Instance_of_instance
			query     = r.FormValue("") // <-- query recuperation
			all_table = []string{"categorie", "email_owner", "post", "user"}
		)
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

			if Check_username == true && Check_email == true { // <-- verify if the user we want to create do not exist

				User.ADD_User_To_BDD(username, pswd_hash, email, rank)
				Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")

			} else { // <-- if wrong selection

				error_message := ""

				if Check_email == false && Check_username == false {
					error_message = "email and username are"
				} else if Check_email == false {
					error_message = "email is"
				} else if Check_username == false {
					error_message = "username is"
				}

				fmt.Fprint(w, "<script> window.alert('This "+error_message+" already in use, try again'); </script>")
				Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
			}

		} else if query == "See_Table" {

			I_I = Database.Select_All_Rows_Table(db, all_table)
			Return_With_Value_Admin(w, r, I_I)

		} else if query == "Backup" {

			I_I = Database.Select_All_Rows_Table(db, all_table)
			Set_Backup(I_I)
			fmt.Fprint(w, "<script>alert('Backup create Succesfully')</script>")
			Config.Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")

		} else if query == "Manage_Users" {
			var (
				table = []string{"user"}
			)
			I_I = Database.Select_All_Rows_Table(db, table)
			Return_With_Value_Admin(w, r, I_I)
		} else if query == "Create_Category" {

			var (
				Create_categorie = r.Form["create_category"][0]
				Check_categorie  = User.Check_If_Exist(Create_categorie, "", "Name", "categorie", "New_categorie")
			)

			if Check_categorie == true {
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
		}

	} else {
		Config.Send_Error(w, r)
		return
	}
}
func Return_With_Value_Admin(w http.ResponseWriter, r *http.Request, I Config.Instance_of_instance) {
	template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/admin/panel_admin.html"))).Execute(w, I)
}

func Set_Backup(I_I Config.Instance_of_instance) {
	var (
		count        = 0
		number_files = strconv.Itoa(count)
		files, err   = ioutil.ReadDir("../backup/")
		t            = time.Now()
		date         = t.Format("2006-01-02")
		backup, _    = json.Marshal(I_I)
	)
	if err != nil {
		log.Fatal(err)
	}
	for range files {
		count++
	}
	if count > 0 {
		number_files = "(" + strconv.Itoa(count) + ")"
	} else {
		number_files = ""
	}

	f, err := os.Create("../backup/Backup_of_" + Config.Bdd.Name + "_" + date + number_files + ".json")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(backup)
}
