package main

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
)

//#------------------------------------------------------------------------------------------------------------# ↓ Page Admin_Panel ↓

func Admin_Panel(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" { // cookie verification to go to admin panel (security)

		s, rank, _ := Check_Cookie(w, r)
		if s == true && rank == "1" {
			Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
		} else { //<-- rank == 4
			Send_Error(w, r)
		}

	} else if r.Method == "POST" { // We accept only methods "POST" and "GET" for security
		var (
			db, err = sql.Open(Bdd.Langage, Bdd.Name)
			I_I     Instance_of_instance
			query   = r.FormValue("") // <-- query recuperation
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
				Check_username = Check_If_Exist(username, "", "Name", "user", "Register")
				Check_email    = Check_If_Exist(email, "", "Email", "user", "Register")
				pswd_hash      = initHashPswd(pswd)
			)

			if Check_username == true && Check_email == true { // <-- verify if the user we want to create do not exist

				ADD_User_To_BDD(username, pswd_hash, email, rank)
				Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")

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
				Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")
			}

		} else if query == "See_Table" {

			I_I = Select_All_Rows_Table(db)
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/admin/panel_admin.html"))).Execute(w, I_I)

		} else if query == "Backup" {

			I_I = Select_All_Rows_Table(db)
			Set_Backup(I_I)
			Return_To_Page(w, r, "../static/templates/admin/panel_admin.html")

		} else if query == "Manage_Users" {

		}

	} else {
		Send_Error(w, r)
		return
	}
}

func Set_Backup(I_I Instance_of_instance) {
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

	f, err := os.Create("../backup/Backup_of_" + Bdd.Name + "_" + date + number_files + ".json")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(backup)
}
