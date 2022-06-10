package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Connected_Status struct {
	User           string
	User_Hased     string
	Rank_Id        string
	Rank_Id_Hashed string
}

var Connected Connected_Status
var templatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return to [Select_Page] ↓

//Return to page Selected need Path (string)
func Return_To_Page(w http.ResponseWriter, r *http.Request, Path string) {
	template.Must(template.ParseFiles(filepath.Join(templatesDir, Path))).Execute(w, " ")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	Return_To_Page(w, r, "../static/templates/managed_pages/404.html")
	// http.Error(w, "Method is not supported.", http.StatusBadRequest) //<-- Print [error] Method is not supported
	// fmt.Fprint(w, http.StatusBadRequest)

}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓

/*
Exemple:
UPDATE user
SET Name = 'test'
WHERE Name = 'New_test';
*/
//Change value on table
func Update_Field(Table, field_table, field_table_two, Last_input, New_input string) {
	var (
		db, err = sql.Open(Bdd.Langage, Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("UPDATE " + Table + " SET " + field_table + " = '" + New_input + "' WHERE " + field_table_two + " = '" + Last_input + "';")
}
func Reset_Password(New_password, Last_password string) {
	New_password = initHashPswd(New_password)
	Update_Field("user", "Pswd", "Pswd", Last_password, New_password)
}

//Server Http
func httpServ() {

	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)
	http.HandleFunc("/validation_mail", Validation_URLbyMail)
	http.HandleFunc("/resend_mail", Resend_Mail)
	http.HandleFunc("/admin_panel", Admin_Panel)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/valide_password", valide_password_page)
	http.HandleFunc("/reset_password_page", reset_password_page)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
