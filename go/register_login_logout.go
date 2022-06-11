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

	emailverifier "github.com/AfterShip/email-verifier"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Logout ↓

//Logout
func logout(w http.ResponseWriter, r *http.Request) {
	del(w, r)
	fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Login ↓

//Login Page
func login(w http.ResponseWriter, r *http.Request) {
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

			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/login.html"))).Execute(w, pos)

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
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
				//--> redirect to index

				//--> redirect to index.html
			} else { // <-- Send Error

				fmt.Fprint(w, "<script> window.alert('Bad password or bad identification, try again.'); </script>")
				template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/login.html"))).Execute(w, pos)

			}

		} else { // <-- If r.Method != Get/Post

			Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Register ↓

//Register Page
func register(w http.ResponseWriter, r *http.Request) {
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
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/register.html"))).Execute(w, pos)

		} else if r.Method == "POST" {
			//<-- check temp bdd
			var (
				query = r.FormValue("") // <-- query recuperation

			)
			if query == "register" {
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
					Email_test_by_dns     = Email_Validation(Email_Register) // <-- Test mail by dns
					Check_temp_user       = Check_If_Exist(User_Register, "", "Name", "temp_user", "Register")
					Check_temp_email      = Check_If_Exist(Email_Register, "", "Email", "temp_user", "Register")
				)

				if Email_test_by_dns == false {
					fmt.Fprint(w, "<script> window.alert('Wrong email "+Email_Register+"'); </script>")
					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/register.html"))).Execute(w, pos)
					return
				} else if Check_temp_user == false || Check_temp_email == false { // <-- Check if exist in temp_user table
					var (
						Error = ""
					)

					if Check_temp_user == false && Check_temp_email == false {
						Error = "Username and email"
					} else if Check_temp_email == false {
						Error = "Email"
					} else if Check_temp_user == false {
						Error = "Username"
					}

					fmt.Fprint(w, "<script> window.alert('"+Error+" alerady used, check your emails or try tomorrow.'); </script>")
					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/register.html"))).Execute(w, pos)
					return
				}

				if Check_User == true && Check_Email == true && Pswd_Register == Pswd_Register_Confirm { // <-- If all is ok

					Register_Smtp(Email_Register, User_Register, user_hash)
					ADD_User_To_Temp(User_Register, Hash_Pswd, Email_Register, user_hash)
					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/register.html"))).Execute(w, pos)

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

					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/register.html"))).Execute(w, pos)
				}

			}

		} else { // <-- If r.Method != Get/Post

			Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Reset password Page ↓
func reset_password_page(w http.ResponseWriter, r *http.Request) {
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

			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)

		} else if r.Method == "POST" {

			r.ParseForm()

			var (
				email_to_reset = r.Form["reset_email"][0]
				email_exist    = Check_If_Exist(email_to_reset, "", "Email", "user", "Register") // true -> don't exist
			)

			if email_exist == true {

				template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
				fmt.Fprint(w, "<script> window.alert('Mail don't exist'); </script>")

			} else {

				var (
					email_tab       = []string{email_to_reset}
					validation_hash = Return_From_Table(email_to_reset, "user", "Reset_password")
				)

				if validation_hash == "error" || len(validation_hash) == 0 {

					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
					fmt.Fprint(w, "<script> window.alert('Contact administrator'); </script>")

				} else {

					Init_Smtp(email_tab, "", validation_hash, "Reset")
					template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/reset_password_page.html"))).Execute(w, pos)
					fmt.Fprint(w, "<script> window.alert('Sent ! Now check your emails.'); </script>")
				}

			}

		} else {

			Send_Error(w, r)

			return
		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Reset password Validation by query ↓
func valide_password_page(w http.ResponseWriter, r *http.Request) {
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
			Return_To_Page(w, r, "../static/templates/managed_pages/valide_password_page.html")
		} else if r.Method == "POST" {

			var (
				query                 = r.FormValue("")
				reset_password        = r.Form["reset_password"][0]
				confim_reset_password = r.Form["confirm_reset_password"][0]
				ok                    = Check_If_Exist(query, "", "Pswd", "user", "validation")
			)

			r.ParseForm()
			if ok == false {
				if reset_password == confim_reset_password {
					Reset_Password(confim_reset_password, query)
					fmt.Fprint(w, `<script> window.alert('Password reset successfully') </script>`)
				} else {
					fmt.Fprint(w, `<script> window.alert('both passwords must be similar') </script>`)
					Return_To_Page(w, r, "../static/templates/managed_pages/valide_password_page.html")
				}

			} else {
				Send_Error(w, r)
				return
			}
		} else {
			Send_Error(w, r)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			return
		}
	} else {
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
	} else if Who_Want == "validation" || Who_Want == "New_categorie" || Who_Want == "temp_user" || Who_Want == "Reset" {
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

//#------------------------------------------------------------------------------------------------------------# ↓ Check func for reset password ↓
func Reset_Password(New_password, Last_password string) {
	New_password = initHashPswd(New_password)
	Update_Field("user", "Pswd", "Pswd", Last_password, New_password)
}

func Return_From_Table(input, table_name, who_want string) string {

	var (
		db, err = sql.Open(Bdd.Langage, Bdd.Name)
		rows    = Select_All_From_DB(db, table_name)
		Rows    []string
		u       = all_bd{}
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
			Rows = append(Rows, strconv.Itoa(*&u.User.Id), *&u.User.Desc, *&u.User.Email, *&u.User.Name, *&u.User.Profile_Picture, *&u.User.Profile_Picture, *&u.User.Pswd, strconv.Itoa(*&u.User.Rank_id))
		} else if who_want == "Email_profil" {
			err := rows.Scan(&u.User.Id, &u.User.Desc, &u.User.Email, &u.User.Name, &u.User.Profile_Picture, &u.User.Pswd, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, strconv.Itoa(*&u.User.Id), *&u.User.Desc, *&u.User.Email, *&u.User.Name, *&u.User.Profile_Picture, *&u.User.Profile_Picture, *&u.User.Pswd, strconv.Itoa(*&u.User.Rank_id))
		}
	}
	input = strings.ToLower(input)
	for i := range Rows {
		Rows[i] = strings.ToLower(Rows[i])
		if Rows[i] == input {
			check = true
		}
		if Rows[i] == input {

			if check == true {

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

//#------------------------------------------------------------------------------------------------------------# ↓ Pages Manage smtp / Temp users / validation Query ↓
func Validation_URLbyMail(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")
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
	if pos.Rank == "4" {
		if r.Method == "GET" {
			Return_To_Page(w, r, "../static/templates/managed_pages/Validation_URLbyMail.html")

		} else if r.Method == "POST" {
			Check_Validation_QueryURL(w, r, query)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)

		} else { // <-- If r.Method != Get/Post

			Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

func Check_Validation_QueryURL(w http.ResponseWriter, r *http.Request, query string) bool {
	var (
		test    = Check_If_Exist(query, "", "validation", "temp_user", "validation")
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

//-------------------------------------------------
func need_validation(email, user string) string {
	var (
		db, err = sql.Open(Bdd.Langage, Bdd.Name)
		rows    = Select_All_From_DB(db, "temp_user")
		Rows    []string
		u       = all_bd{}
		check   = false
	)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.validation)
		if err != nil {
			log.Fatal(err)
		}
		Rows = append(Rows, strconv.Itoa(*&u.Temp_user.Id), *&u.Temp_user.Name, *&u.Temp_user.Email, *&u.Temp_user.Pswd, *&u.Temp_user.validation)

	}

	for i := range Rows {
		if Rows[i] == user {
			check = true
		}
		if Rows[i] == email {
			if check == true {
				return Rows[i+2]
			}
		}
	}
	return "error"
}

//--------------------------------------------------
func Resend_Mail(w http.ResponseWriter, r *http.Request) {
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
			template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/resend_mail.html"))).Execute(w, pos)
		} else if r.Method == "POST" {

			r.ParseForm()
			var (
				User_Resend  = r.Form["user_resend"][0]
				Email_Resend = r.Form["email_resend"][0]
				hash         = need_validation(Email_Resend, User_Resend)
			)

			if hash == "error" {
				fmt.Fprint(w, "<script> window.alert('Bad emailor username, try again'); </script>")
				template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/templates/resend_mail.html"))).Execute(w, pos)
			} else {
				Register_Smtp(Email_Resend, User_Resend, hash)
				fmt.Fprint(w, "<script> window.alert('Mail sent successfully'); </script>")
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
			}

		} else { // <-- If r.Method != Get/Post

			Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Check mail(register) By DNS ↓
var (
	verifier = emailverifier.NewVerifier()
)

// Email verifier by DNS
func Email_Validation(email string) bool {
	_, err := verifier.Verify(email)
	if err != nil {
		return false
	}
	return true
}
