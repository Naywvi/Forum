package user

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	Config "forum/config"
	Database "forum/database"
)

func Register(w http.ResponseWriter, r *http.Request) { //Register Page
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
			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/register.html"))).Execute(w, pos)

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
					Hash_Pswd             = Database.InitHashPswd(Pswd_Register)
					user_hash             = Database.InitHashPswd(User_Register)
					Email_test_by_dns     = Email_Validation(Email_Register) // <-- Test mail by dns
					Check_temp_user       = Check_If_Exist(User_Register, "", "Name", "temp_user", "Register")
					Check_temp_email      = Check_If_Exist(Email_Register, "", "Email", "temp_user", "Register")
				)

				if Email_test_by_dns == false {
					fmt.Fprint(w, "<script> window.alert('Wrong email "+Email_Register+"'); </script>")
					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/register.html"))).Execute(w, pos)
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
					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/register.html"))).Execute(w, pos)
					return
				}

				if Check_User == true && Check_Email == true && Pswd_Register == Pswd_Register_Confirm { // <-- If all is ok

					Register_Smtp(Email_Register, User_Register, user_hash)
					ADD_User_To_Temp(User_Register, Hash_Pswd, Email_Register, user_hash)
					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/register.html"))).Execute(w, pos)

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

					template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/register.html"))).Execute(w, pos)
				}

			}

		} else { // <-- If r.Method != Get/Post

			Config.Send_Error(w, r)
			return

		}
	} else {
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
		return
	}
}
