package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	Config "forum/config"
	Database "forum/database"
)

func Valide_password_page(w http.ResponseWriter, r *http.Request) {
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
	if pos.Rank != "4" {
		if r.Method == "GET" {
			Config.Return_To_Page(w, r, "../static/templates/managed_pages/valide_password_page.html")
		} else if r.Method == "POST" {

			var (
				query                 = r.FormValue("")
				reset_password        = r.Form["reset_password"][0]
				confim_reset_password = r.Form["confirm_reset_password"][0]
				ok                    = Check_If_Exist(query, "", "Pswd", "user", "validation")
			)

			r.ParseForm()
			if !ok {
				if reset_password == confim_reset_password {
					Reset_Password(confim_reset_password, query)
					fmt.Fprint(w, `<script> window.alert('Password reset successfully') </script>`)
				} else {
					fmt.Fprint(w, `<script> window.alert('both passwords must be similar') </script>`)
					Config.Return_To_Page(w, r, "../static/templates/managed_pages/valide_password_page.html")
				}

			} else {
				Config.Send_Error(w, r)
				return
			}
		} else {
			Config.Send_Error(w, r)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			return
		}
	} else {
		return
	}
}

//Multiple func Who_Whant (--> "Register" or "Login")
//Select field for log/register in bdd
func Check_If_Exist(input, input2, check_field, In_This_Table, Who_Want string) bool {

	var (
		db, _ = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		Rows  []string
		rows  *sql.Rows
		I     = Config.Instance_Bdd{}
		u     = Config.All_bd{}
		index = &u.User.Name //<-- Default as Name & redefine on if
	)

	if Who_Want == "Register" { //<-- Define var by Who_want
		rows = Database.Select_Field_From_DB(db, check_field, In_This_Table)
		input = strings.ToLower(input)
	} else if Who_Want == "Login" {
		rows = Database.Select_All_From_DB(db, "user")
	} else if Who_Want == "validation" || Who_Want == "New_categorie" || Who_Want == "temp_user" || Who_Want == "Reset" {
		rows = Database.Select_Field_From_DB(db, check_field, In_This_Table)
	}
	for rows.Next() {

		if Who_Want == "Register" { //<-- Select field
			if check_field == "Email" {
				index = &u.User.Email
			}
			err := rows.Scan(index)
			if err != nil {
				log.Fatal(err)
			}

			Rows = append(Rows, *index)

		} else if Who_Want == "Login" {

			err := rows.Scan(&u.User.Id, &u.User.Name, &u.User.Pswd, &u.User.Desc, &u.User.Email, &u.User.Profile_Picture, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}

			u.User.Name = strings.ToLower(u.User.Name) //<-- Check in lower case to be sure
			u.User.Email = strings.ToLower(u.User.Email)
			I.I = append(I.I, u) //<-- Instance of User{}
		} else if Who_Want == "validation" {

			err := rows.Scan(index)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, *index)

		} else if Who_Want == "New_categorie" {

			err := rows.Scan(&u.Categorie.Name)
			if err != nil {
				log.Fatal(err)
			}

			Rows = append(Rows, u.Categorie.Name)
		} else if Who_Want == "temp_user" {

			if check_field == "Email" {
				index = &u.Temp_user.Email
			} else if check_field == "Name" {
				index = &u.Categorie.Name
			}

			err := rows.Scan(&u.Temp_user.Name, &u.Temp_user.Email)
			if err != nil {
				log.Fatal(err)
			}

			Rows = append(Rows, *index)
		} else if Who_Want == "Reset" {
			index = &u.User.Email
			err := rows.Scan(index)
			if err != nil {
				log.Fatal(err)
			}

			Rows = append(Rows, *index)

		}
	}

	if Who_Want == "Register" || Who_Want == "validation" || Who_Want == "New_categorie" || Who_Want == "temp_user" || Who_Want == "Reset" { //<-- Select return
		return Check_Login_Or_Register(nil, "", "", Who_Want, input, Rows)
	} else if Who_Want == "Login" {
		return Check_Login_Or_Register(&I, input, input2, Who_Want, "", nil)
	}

	return false
}

func Return_From_Table(input, table_name, who_want string) string {

	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		rows    = Database.Select_All_From_DB(db, table_name)
		Rows    []string
		u       = Config.All_bd{}
		check   = false
	)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if who_want == "Reset_password" {
			err := rows.Scan(&u.User.Id, &u.User.Desc, &u.User.Email, &u.User.Name, &u.User.Profile_Picture, &u.User.Pswd, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, strconv.Itoa(u.User.Id), u.User.Desc, u.User.Email, u.User.Name, u.User.Profile_Picture, u.User.Profile_Picture, u.User.Pswd, strconv.Itoa(u.User.Rank_id))
		} else if who_want == "Email_profil" {
			err := rows.Scan(&u.User.Id, &u.User.Desc, &u.User.Email, &u.User.Name, &u.User.Profile_Picture, &u.User.Pswd, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, strconv.Itoa(u.User.Id), u.User.Desc, u.User.Email, u.User.Name, u.User.Profile_Picture, u.User.Profile_Picture, u.User.Pswd, strconv.Itoa(u.User.Rank_id))
		}
	}
	input = strings.ToLower(input)
	for i := range Rows {
		Rows[i] = strings.ToLower(Rows[i])
		if Rows[i] == input {
			check = true
		}
		if Rows[i] == input {

			if check {

				if who_want == "Reset_password" {
					return Rows[i-2]
				} else if who_want == "Email_profil" { // savoir si faut mettre +3 ou -1943 depuis bdd ou juste au dessus ordre des variables?
					return Rows[i+3]
				}

			}
		}

	}

	return "error"
}

//#------------------------------------------------------------------------------------------------------------# ↓ Check func for register or login ↓

//Check in dbb if exist
func Check_Login_Or_Register(I *Config.Instance_Bdd, identifier, pswd, Who_whant, input string, Rows []string) bool {

	if Who_whant == "Login" {

		for _, i := range I.I { //<-- Check in I (= instance of table)
			if identifier == i.User.Email || identifier == i.User.Name {
				//<<<< Set for personnal Cookie
				Config.Connected.User = i.User.Name
				Config.Connected.User_Hased = Encrypt_Cookie(Config.Connected.User)
				Config.Connected.Rank_Id = strconv.Itoa(i.User.Rank_id)
				Config.Connected.Rank_Id_Hashed = Encrypt_Cookie(strconv.Itoa(i.User.Rank_id))
				//<<<< Set for personnal Cookie
				return Database.CheckPasswordHash(pswd, i.User.Pswd)
			}
		}
		return false

	} else if Who_whant == "Register" {

		for i := range Rows {
			low := strings.ToLower(Rows[i]) //<-- Check in lower case to be sure
			if low == input {
				return false
			}
		}

		return true
	} else if Who_whant == "validation" || Who_whant == "New_categorie" || Who_whant == "temp_user" || Who_whant == "Reset" {

		for i := range Rows {
			if Rows[i] == input {
				return false
			}
		}

		return true
	}
	return false
}
