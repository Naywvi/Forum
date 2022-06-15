package user

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

	"github.com/k3a/html2text"

	Config "forum/config"
	Database "forum/database"
)

func Delete_Account(user string) {
	db, err := sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("DELETE FROM profil WHERE user = '" + user + "'")
	db.Exec("DELETE FROM user WHERE name = '" + user + "'")
}

func Delete_Category(category string) {
	db, err := sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok", category)
	db.Exec("DELETE FROM post WHERE id_cat = '" + category + "'")
	db.Exec("DELETE FROM categorie WHERE name = '" + category + "'")
}

func Profildeleted(w http.ResponseWriter, r *http.Request) {
	var (
		_, statement, user = Check_Cookie(w, r)
	)
	query := r.FormValue("User")
	//<<< --- Check rank
	if statement != "4" {

		if r.Method == "GET" {
			fmt.Println(query)
			if (len(query) > 0 && statement == "2") || (len(query) > 0 && statement == "1") {
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			} else {
				Delete_Account(user)
				Logout(w, r)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Edit desc ↓

func Edit_desc(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	query_other_desc := r.FormValue("User")
	query_other_send := r.FormValue("send")

	type Statement_of_user struct {
		User                string
		Rank                string
		Desc                string
		Descedit            string
		User_connected      string
		User_connected_rank string
		Subject_submited    string
	}

	//<<< --- Check rank

	var (
		_, statement, User = Check_Cookie(w, r)
		pos                = Statement_of_user{}
	)

	//<<< --- Check rank
	if statement != "4" {

		if r.Method == "GET" {

			var (
				Check_user_edit_desc = Check_If_Exist(query_other_desc, "", "name", "user", "Register")
				result               *sql.Rows
				result_profil        []string
			)

			if len(query_other_desc) > 0 {

				if Check_user_edit_desc == false {

					result = Database.Select_column("profil", "user", query_other_desc) //Rows

				} else {
					Config.Send_Error(w, r)

					return
				}

			} else {

				result = Database.Select_column("profil", "user", User) //Rows

			}

			result_profil = Return_Profil(result)
			//<<<<
			pos.Descedit = html2text.HTML2Text(result_profil[6])
			pos.Desc = result_profil[6]
			pos.Rank = result_profil[7]
			pos.User = result_profil[1]
			pos.User_connected = User
			pos.User_connected_rank = statement
			pos.Subject_submited = result_profil[4]

			//<<<<

			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/managed_pages/edit_desc_profile.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

			desc_edit := r.Form["description"][0]

			if query == "send" {

				if len(desc_edit) > 2000 || len(desc_edit) == 0 {
					fmt.Fprint(w, "<script> window.alert('Description too long.'); </script>")
					fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil/edit"; </script>`)
					return
				}

				Database.Update_Field("profil", "Desc", desc_edit, "user", User)

				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil?=`+User+`"; </script>`)

			} else if len(query_other_send) > 0 {

				if len(desc_edit) > 2000 || len(desc_edit) == 0 {
					fmt.Fprint(w, "<script> window.alert('Description too long.'); </script>")
					fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil/edit"; </script>`)
					return
				}

				Database.Update_Field("profil", "Desc", desc_edit, "user", query_other_send)

				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/profil?=`+User+`"; </script>`)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Init profil ↓
//Profil Page
func Return_Profil(rows *sql.Rows) []string {
	var (
		test   Config.All_bd
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

func Profil(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
	type Statement_of_user struct {
		User                string
		Rank                string
		Email               string
		Joined              string
		Last_time_connected string
		Subject_submit      string
		Desc                string
		User_connected      string
		User_connected_rank string
	}
	//<<< --- Check rank

	var (
		_, statement, User_connected = Check_Cookie(w, r)
		pos                          = Statement_of_user{}
	)

	//<<< --- Check rank
	if statement != "4" {

		if r.Method == "GET" {

			Check_user := Check_If_Exist(query, "", "name", "user", "Register")

			if Check_user == false {
				var (
					result        = Database.Select_column("profil", "user", query) //Rows
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
				pos.User_connected = User_connected
				pos.User_connected_rank = statement

				//<<<<
				template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/profil.html"))).Execute(w, pos)
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
		return
	}

}

//Create Default profil
func New_Profil(User, Email, rank_id string) { //User,Joined,Last_time_connected,Subjet_submit,Email,Desc
	var (
		db, err            = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		joined             = time.Now()
		Default_profil_arr = []string{"'" + User + "','" + joined.String() + "','" + joined.String() + "','0','" + Email + "','none_desc','" + rank_id + "'"}
		Default_profil     = strings.Join(Default_profil_arr, "")
	)
	if err != nil {
		log.Fatal(err)
	}
	Database.Inser_In_To_DB(db, Default_profil, "profil", Database.Extract_File("../bdd/profil_table.sql", 11, 12))
}

//#------------------------------------------------------------------------------------------------------------# ↓ Add desc ↓
