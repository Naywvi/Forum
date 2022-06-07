package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
func Del_User_From_Table(db *sql.DB, rows *sql.Rows, table, name_deleted, who_want string) { //all time send i of deleter
	var (
		Rows  []string
		u     = all_bd{}
		marge = 0
		id    = ""
		index = 0
	)

	for rows.Next() {
		marge = 4 // <-- De combien je recule pour avoir l'id dans la table afin de le delect (vérification par la "validation field")
		if who_want == "validation" {
			err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.validation)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, strconv.Itoa(*&u.Temp_user.Id), *&u.Temp_user.Name, *&u.Temp_user.Email, *&u.Temp_user.Pswd, *&u.Temp_user.validation)
		}

	}

	for i := range Rows {

		if Rows[i] == name_deleted {
			index = i
			id = Rows[i-marge]
			break
		}

	}
	db.Exec("DELETE FROM " + table + " WHERE id = " + id)
	fmt.Println("Validation done successfully", index)
	// if who_want == "validation" {
	// 	ADD_User_To_BDD(Rows[index-3], Rows[index-1], Rows[index-2], "3")
	// }
}
func Check_Validation_QueryURL(w http.ResponseWriter, r *http.Request, query string) bool {
	test := Check_If_Exist(query, "", "validation", "temp_user", "validation")
	var (
		db, err = sql.Open(Bdd.Langage, Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}
	if test == false {

		Del_User_From_Table(db, Select_All_From_DB(db, "temp_user"), "temp_user", query, "validation")
		print("coucou")
	}
	print("salut")
	Send_Error(w, r)
	return false
}

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	Return_To_Page(w, r, "../static/templates/managed_pages/404.html")
	// http.Error(w, "Method is not supported.", http.StatusBadRequest) //<-- Print [error] Method is not supported
	// fmt.Fprint(w, http.StatusBadRequest)

}
func Validation_URLbyMail(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	if r.Method == "GET" {
		Return_To_Page(w, r, "../static/templates/managed_pages/Validation_URLbyMail.html")

	} else if r.Method == "POST" {
		Check_Validation_QueryURL(w, r, query)
	} else { // <-- If r.Method != Get/Post

		Send_Error(w, r)
		return

	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Selection & init http_serv ↓

//Server Http
func httpServ() {

	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)
	http.HandleFunc("/validation_mail", Validation_URLbyMail)

	http.HandleFunc("/admin_panel", Admin_Panel)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	fmt.Println("Started https serv successfully on http://localhost:1010")
	http.ListenAndServe(":1010", nil)

}
