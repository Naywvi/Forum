package user

import (
	"database/sql"
	"fmt"
	Config "forum/config"
	Database "forum/database"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	emailverifier "github.com/AfterShip/email-verifier"
)

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
			Config.Return_To_Page(w, r, "../static/templates/managed_pages/Validation_URLbyMail.html")

		} else if r.Method == "POST" {
			Check_Validation_QueryURL(w, r, query)
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)

		} else { // <-- If r.Method != Get/Post

			Config.Send_Error(w, r)
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
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	)

	if err != nil {
		log.Fatal(err)
	}

	if test {
		Del_User_From_Table(db, Database.Select_All_From_DB(db, "temp_user"), "temp_user", query, "validation")
	}

	Config.Send_Error(w, r)
	return false
}

//-------------------------------------------------
func need_validation(email, user string) string {
	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		rows    = Database.Select_All_From_DB(db, "temp_user")
		Rows    []string
		u       = Config.All_bd{}
		check   = false
	)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {

		err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.Validation)
		if err != nil {
			log.Fatal(err)
		}
		Rows = append(Rows, strconv.Itoa(*&u.Temp_user.Id), *&u.Temp_user.Name, *&u.Temp_user.Email, *&u.Temp_user.Pswd, *&u.Temp_user.Validation)
	}

	for i := range Rows {
		if Rows[i] == user {
			check = true
		}
		if Rows[i] == email {
			if check {
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
		fmt.Print(r.Method)
		if r.Method == "GET" {
			template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/resend_mail.html"))).Execute(w, pos)

		} else if r.Method == "POST" {
			fmt.Print("RENTRE")
			r.ParseForm()
			var (
				User_Resend  = r.Form["user_resend"][0]
				Email_Resend = r.Form["email_resend"][0]
				hash         = need_validation(Email_Resend, User_Resend)
			)

			if hash == "error" {
				fmt.Fprint(w, "<script> window.alert('Bad emailor username, try again'); </script>")
				template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/resend_mail.html"))).Execute(w, pos)
			} else {
				Register_Smtp(Email_Resend, User_Resend, hash)
				fmt.Fprint(w, "<script> window.alert('Mail sent successfully'); </script>")
				fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Check mail(register) By DNS ↓
var (
	verifier = emailverifier.NewVerifier()
)

// Email verifier by DNS
func Email_Validation(email string) bool {
	_, err := verifier.Verify(email)
	return err == nil
}
