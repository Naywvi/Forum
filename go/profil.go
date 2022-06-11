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

//#------------------------------------------------------------------------------------------------------------# ↓ Init profil ↓
//Profil Page
func Return_Profil(rows *sql.Rows) []string {
	var (
		test   all_bd
		result []string
	)
	for rows.Next() {
		err := rows.Scan(&test.Profil.Id, &test.Profil.User, &test.Profil.Joined, &test.Profil.Last_time_connected, &test.Profil.Subjet_submit, &test.Profil.Email, &test.Profil.Desc, &test.Profil.Rank_id_profil)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, strconv.Itoa(test.Profil.Id), test.Profil.User, test.Profil.Joined, test.Profil.Last_time_connected, test.Profil.Subjet_submit, test.Profil.Email, test.Profil.Desc, test.Profil.Rank_id_profil)
	}
	return result
}
func profil(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Statement_of_user struct {
		User                string
		Rank                string
		Email               string
		Joined              string
		Last_time_connected string
		Subject_submit      string
		Desc                string
	}
	//<<< --- Check rank

	var (
		_, statement, _ = Check_Cookie(w, r)
		pos             = Statement_of_user{}
	)

	//<<< --- Check rank
	if statement != "4" {

		if r.Method == "GET" {
			var (
				result        = Select_column("profil", "user", query) //Rows
				result_profil = Return_Profil(result)
			)
			//<<<<
			pos.User = result_profil[1]
			pos.Joined = result_profil[2][0:10]
			pos.Last_time_connected = result_profil[3][0:10]
			pos.Subject_submit = result_profil[4]
			pos.Email = result_profil[5]
			pos.Desc = result_profil[6]
			pos.Rank = result_profil[7]
			//<<<<
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/profil.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

		} else {

			Send_Error(w, r)

			return
		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//Create Default profil
func New_Profil(User, Email string) { //User,Joined,Last_time_connected,Subjet_submit,Email,Desc
	var (
		db, err            = sql.Open(Bdd.Langage, Bdd.Name)
		joined             = time.Now()
		Default_profil_arr = []string{"'" + User + "','" + joined.String() + "','" + joined.String() + "','0','" + Email + "','none_desc','3'"}
		Default_profil     = strings.Join(Default_profil_arr, "")
	)
	if err != nil {
		log.Fatal(err)
	}
	Inser_In_To_DB(db, Default_profil, "profil", Extract_File("../bdd/profil_table.sql", 11, 12))
}

//#------------------------------------------------------------------------------------------------------------# ↓ Add desc ↓
