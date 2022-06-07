package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Logout ↓

//Logout
func logout(w http.ResponseWriter, r *http.Request) {
	del(w, r)
	fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Login ↓

//Login Page
func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		Return_To_Page(w, r, "../static/templates/login.html")

	} else if r.Method == "POST" {

		r.ParseForm()

		var (
			mail_login = r.Form["mail_login"][0]
			pswd_login = r.Form["password_login"][0]
			Check      = Check_If_Exist(mail_login, pswd_login, "Pswd", "user", "Login")
		)

		if Check == true {
			SettCookie(w, r) //send cookie first
			fmt.Fprint(w, `<script> window.alert('Your are connected') </script>`)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
			//--> redirect to index

			//--> redirect to index.html
		} else { // <-- Send Error

			fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
			Return_To_Page(w, r, "../static/templates/login.html")

		}

	} else { // <-- If r.Method != Get/Post

		Send_Error(w, r)
		return

	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Register ↓

//Register Page
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		Return_To_Page(w, r, "../static/templates/register.html")

	} else if r.Method == "POST" {
		//<-- check temp bdd
		r.ParseForm()

		var (
			User_Register         = r.Form["username_register"][0]
			Email_Register        = r.Form["email_register"][0]
			Pswd_Register         = r.Form["password_register"][0]
			Pswd_Register_Confirm = r.Form["password_register_confirm"][0]
			Check_User            = Check_If_Exist(User_Register, "", "Name", "user", "Register")
			Check_Email           = Check_If_Exist(Email_Register, "", "Email", "user", "Register")
			Hash_Pswd             = initHashPswd(Pswd_Register)
			user_hash             = initHashPswd(User_Register)
		)

		if Check_User == true && Check_Email == true && Pswd_Register == Pswd_Register_Confirm { // <-- If all is ok

			Register_Smtp(Email_Register, User_Register, user_hash)
			ADD_User_To_Temp(User_Register, Hash_Pswd, Email_Register, user_hash)
			Return_To_Page(w, r, "../static/templates/managed_pages/after_register.html")

			//--> redirect to index
			return

		} else { // <-- Check the wrong selection

			error_message := ""
			if Pswd_Register == Pswd_Register_Confirm {
				if Check_Email == false && Check_User == false {
					error_message = "email and username are"
				} else if Check_Email == false {
					error_message = "email is"
				} else if Check_User == false {
					error_message = "username is"
				}
				fmt.Fprint(w, "<script> window.alert('This "+error_message+" already in use, try again'); </script>")
			} else {
				fmt.Fprint(w, "<script> window.alert('Bad password confirmation, try again'); </script>")
			}

			Return_To_Page(w, r, "../static/templates/register.html")
		}

	} else { // <-- If r.Method != Get/Post

		Send_Error(w, r)
		return

	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Select field for register or login ↓

//Multiple func Who_Whant (--> "Register" or "Login")
//Select field for log/register in bdd
func Check_If_Exist(input, input2, check_field, In_This_Table, Who_Want string) bool {

	var (
		db, _ = sql.Open(Bdd.Langage, Bdd.Name)
		Rows  []string
		rows  *sql.Rows
		I     = Instance_Bdd{}
		u     = all_bd{}
		index = &u.User.Name //<-- Default as Name & redefine on if
	)

	if Who_Want == "Register" { //<-- Define var by Who_want
		rows = Select_Field_From_DB(db, check_field, In_This_Table)
		input = strings.ToLower(input)
	} else if Who_Want == "Login" {
		rows = Select_All_From_DB(db, "user")
	} else if Who_Want == "validation" {
		rows = Select_Field_From_DB(db, check_field, In_This_Table)
	} else if Who_Want == "New_categorie" {
		rows = Select_Field_From_DB(db, check_field, In_This_Table)
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
		}

	}

	if Who_Want == "Register" || Who_Want == "validation" || Who_Want == "New_categorie" { //<-- Select return
		return Check_Login_Or_Register(nil, "", "", Who_Want, input, Rows)
	} else if Who_Want == "Login" {
		return Check_Login_Or_Register(&I, input, input2, Who_Want, "", nil)
	}

	return false
}

//#------------------------------------------------------------------------------------------------------------# ↓ Check func for register or login ↓

//Check in dbb if exist
func Check_Login_Or_Register(I *Instance_Bdd, identifier, pswd, Who_whant, input string, Rows []string) bool {

	if Who_whant == "Login" {

		for _, i := range I.I { //<-- Check in I (= instance of table)
			if identifier == i.User.Email || identifier == i.User.Name {
				//<<<< Set for personnal Cookie
				Connected.User = i.User.Name
				Connected.User_Hased = Encrypt_Cookie(Connected.User)
				Connected.Rank_Id = strconv.Itoa(i.User.Rank_id)
				Connected.Rank_Id_Hashed = Encrypt_Cookie(strconv.Itoa(i.User.Rank_id))
				//<<<< Set for personnal Cookie
				return CheckPasswordHash(pswd, i.User.Pswd)
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
	} else if Who_whant == "validation" || Who_whant == "New_categorie" {

		for i := range Rows {
			if Rows[i] == input {
				return false
			}
		}

		return true
	}
	return false
}

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Manage smtp / Temp users / validation Query ↓
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
	if who_want == "validation" {
		ADD_User_To_BDD(Rows[index-3], Rows[index-1], Rows[index-2], "3")
	}
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

	}

	Send_Error(w, r)
	return false
}
