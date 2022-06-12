package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Extract sql-file & return it (select interval in file with end/ start)
func Extract_File(file_sql string, start, end int) string {
	var (
		text    = ""
		count   = 0
		file, _ = os.Open(file_sql)
		scanner = bufio.NewScanner(file)
	)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if count == end {
			return text
		}
		if count >= start-1 {
			text += scanner.Text() + "\n"
		}
		count++
	}
	return text
}

//#------------------------------------------------------------------------------------------------------------# ↓ Add user to table ↓

func Print_Rows(rows *sql.Rows, table string) []all_bd {
	var (
		I Instance_Bdd
		u = all_bd{}
	)

	for rows.Next() {

		if table == "user" {
			err := rows.Scan(&u.User.Id, &u.User.Name, &u.User.Pswd, &u.User.Desc, &u.User.Email, &u.User.Profile_Picture, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}

		} else if table == "email_owner" {
			err := rows.Scan(&u.Smtp.Email, &u.Smtp.Pass)
			if err != nil {
				log.Fatal(err)
			}
		} else if table == "categorie" {
			err := rows.Scan(&u.Categorie.Id, &u.Categorie.Name)
			if err != nil {
				log.Fatal(err)
			}
		} else if table == "post" {
			err := rows.Scan(&u.Post.Id, &u.Post.Id_cat, &u.Post.Title_post, &u.Post.Content, &u.Post.Likes, &u.Post.Posted_user, &u.Post.Nb_Reply, &u.Post.Last_Posted)
			if err != nil {
				log.Fatal(err)
			}
		} else if table == "temp_user" {
			err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.validation)
			if err != nil {
				log.Fatal(err)
			}
		} else if table == "profil" {
			err := rows.Scan(&u.Profil.Id, &u.Profil.User, &u.Profil.Joined, &u.Profil.Last_time_connected, &u.Profil.Subjet_submit, &u.Profil.Email, &u.Profil.Desc, &u.Profil.Rank_id_profil)
			if err != nil {
				log.Fatal(err)
			}
		}
		I.I = append(I.I, u)

	}
	return I.I

}

/*
SELECT
	*
FROM
	user
WHERE
	id = 1
*/
func Select_column(Table_name, Table_field, input string) *sql.Rows { //only string
	var (
		db, err = sql.Open(Bdd.Langage, Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}
	result, _ := db.Query("SELECT * FROM " + Table_name + " WHERE " + Table_field + " = '" + input + "';")
	return result
}

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

//Del user from table
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
	if who_want == "validation" {
		ADD_User_To_BDD(Rows[index-3], Rows[index-1], Rows[index-2], "3")
	}
}

//Add user to temp_user
func ADD_User_To_Temp(name, pswd, email, user_hash string) {
	var (
		db, err          = sql.Open(Bdd.Langage, Bdd.Name)
		Default_user_arr = []string{"'" + name + "','" + pswd + "','" + email + "','" + user_hash + "'"}
		Default_user     = strings.Join(Default_user_arr, "")
	)
	if err != nil {
		log.Fatal(err)
	}
	Inser_In_To_DB(db, Default_user, "temp_user", Extract_File("../bdd/temp_user_table.sql", 8, 9))
}

//Add default user to bdd
func ADD_User_To_BDD(name, pswd, email, rank_id string) {
	var (
		db, err          = sql.Open(Bdd.Langage, Bdd.Name)
		Default_user_arr = []string{"'" + name + "','" + pswd + "','none_desc','" + email + "','none_picture','" + rank_id + "'"}
		Default_user     = strings.Join(Default_user_arr, "")
	)

	if err != nil {
		log.Fatal(err)
	}

	Inser_In_To_DB(db, Default_user, "user", Extract_File("../bdd/user_table.sql", 11, 12))
	New_Profil(name, email, rank_id)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Select on table ↓

func Select_All_From_DB(db *sql.DB, table string) *sql.Rows {
	result, _ := db.Query("SELECT * FROM " + table)
	return result
}
func Select_Field_From_DB(db *sql.DB, field, table string) *sql.Rows {
	result, _ := db.Query("SELECT " + field + " FROM " + table)
	return result
}
func Select_All_Rows_Table(db *sql.DB, table []string) Instance_of_instance {
	var (
		I   Instance_Bdd
		I_I Instance_of_instance
	)
	for i := range table {
		I.I = Print_Rows(Select_All_From_DB(db, table[i]), table[i])
		I_I.I = append(I_I.I, I) // <-- To send on one template
	}
	return I_I
}

//#------------------------------------------------------------------------------------------------------------# ↓ Init Add to table ↓

//Auto-Create Table
func Init_Database(table_name, txt string) *sql.DB {
	db, err := sql.Open(Bdd.Langage, Bdd.Name)
	if err != nil {
		log.Fatal(err)
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table_name + "(" + txt + ")"
	db.Exec(sql)
	return db
}

//Auto increment field & Value on table
func Inser_In_To_DB(db *sql.DB, var_receive, table_name, table_field string) (int64, error) { //
	result, err := db.Exec("INSERT INTO " + table_name + " (" + table_field + ")" + " VALUES (" + var_receive + ")")
	if err != nil {
		log.Fatal(err)
	}
	return result.LastInsertId()
}

//#------------------------------------------------------------------------------------------------------------# ↓ Init dbb & featers ↓

//Selection input of shell during init bdd
func Terminal_Init_Table(who_want string) string {

	fmt.Println("---------------------")

	if who_want == "add_user_table" {

		fmt.Println("Create a default user to connect on the web site.")
		fmt.Print("username ->")

		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password ->")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)
		fmt.Println("User -> 3 | Moderator -> 2 | Admin -> 1")
		fmt.Print("rank_id ->")

		rank, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		rank = strings.TrimSpace(rank)
		fmt.Println("---------------------")

		return "'" + username + "','" + pswd + "','none_desc','" + mail + "','none_picture','" + rank + "'"

	} else if who_want == "email_verification_table" {

		fmt.Println("Write your mail to send request register.")
		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password -> ")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n') //<-- Email need \n to connect
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		return "'" + mail + "','" + pswd + "'"
	} else if who_want == "Bdd_Name" {
		fmt.Println("Choose name of your Data_Base.")
		Name_Bd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		Bdd.Name = strings.TrimSpace(Name_Bd) + ".db"
		Bdd.Langage = "sqlite3"

		fmt.Println(Bdd.Name, "Is the database selected. Do you validate ? [y] [n] ")

		Check, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		Check = strings.TrimSpace(Check)

		if Check == "n" || Check == "N" || Check == "no" || Check == "NO" {
			return Terminal_Init_Table(who_want)
		} else if Check == "y" || Check == "Y" || Check == "YES" || Check == "yes" {
			fmt.Println("You are ready now")
		} else {
			fmt.Println("Wrong selection")
			fmt.Println("---------------------")
			return Terminal_Init_Table(who_want)
		}

	} else if who_want == "temp_user" {
		fmt.Println("Create a temp user to check the validation table.")
		fmt.Print("username ->")

		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)
		mail_hash := initHashPswd(mail)

		fmt.Print("password ->")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		fmt.Println("---------------------")
		return "'" + username + "','" + pswd + "','" + mail + "','" + mail_hash + "'"
	}
	fmt.Println("---------------------")
	return ""
}

//Simple print on shell
func Is_Ok(Printable, Second_Printable string) {

	fmt.Println("> " + Printable + " Table was successfully created")

	if len(Second_Printable) > 0 {
		fmt.Println("-> " + Second_Printable + " was successfully created\n")
	} else {
		fmt.Println("")
	}

}

//#------------------------------------------------------------------------------------------------------------# ↓ Init db_test to start ↓
func categorie() {
	Inser_In_To_DB(Init_Database("categorie", Extract_File("../bdd/categorie_table.sql", 0, 2)), Extract_File("../bdd/categorie_table.sql", 8, 10), "categorie", Extract_File("../bdd/categorie_table.sql", 5, 6))
	Is_Ok("categorie", "test_categorie")
}

func post() {
	Inser_In_To_DB(Init_Database("post", Extract_File("../bdd/post_table.sql", 0, 8)), Extract_File("../bdd/post_table.sql", 15, 16), "post", Extract_File("../bdd/post_table.sql", 11, 12))
	Is_Ok("Post", "Test_post")
}
func comment() {
	Inser_In_To_DB(Init_Database("comment", Extract_File("../bdd/comment_table.sql", 0, 6)), Extract_File("../bdd/comment_table.sql", 13, 14), "comment", Extract_File("../bdd/comment_table.sql", 9, 10))
	Is_Ok("Comment", "Comment")
}
func user() {
	Inser_In_To_DB(Init_Database("user", Extract_File("../bdd/user_table.sql", 0, 8)), Terminal_Init_Table("add_user_table"), "user", Extract_File("../bdd/user_table.sql", 11, 12))
	Is_Ok("User", "")
}
func temp_user() {
	Inser_In_To_DB(Init_Database("temp_user", Extract_File("../bdd/temp_user_table.sql", 0, 5)), Terminal_Init_Table("temp_user"), "temp_user", Extract_File("../bdd/temp_user_table.sql", 8, 9))
	Is_Ok("temp_user", "")
}

func email_verification() {
	Inser_In_To_DB(Init_Database("email_owner", Extract_File("../bdd/email_verification_table.sql", 0, 2)), Terminal_Init_Table("email_verification_table"), "email_owner", Extract_File("../bdd/email_verification_table.sql", 5, 6)) //<-- Os.Args email verification
}
func profilt() {
	fmt.Println(Extract_File("../bdd/profil_table.sql", 0, 8))
	fmt.Println(Extract_File("../bdd/profil_table.sql", 11, 12))
	fmt.Println(Extract_File("../bdd/profil_table.sql", 15, 16))
	Inser_In_To_DB(Init_Database("profil", Extract_File("../bdd/profil_table.sql", 0, 8)), Extract_File("../bdd/profil_table.sql", 15, 16), "profil", Extract_File("../bdd/profil_table.sql", 11, 12))
	Is_Ok("profil", "")
}

//#------------------------------------------------------------------------------------------------------------# ↓ init bd ↓

//Init bdd
func InitBDD() {
	if _, err := os.Stat("./" + Bdd.Name); err == nil { //<-- If bdd exist
		fmt.Println("The bdd, " + Bdd.Name + " is already here")

	} else if errors.Is(err, os.ErrNotExist) { //<-- If bdd not exist > Re create
		fmt.Println("--->Create a new db<--- ?")
		email_verification()
		temp_user()
		user()
		categorie()
		post()
		comment()
		profilt()
		fmt.Println("Bdd, " + Bdd.Name + " was successfully created, you are ready :)\n")
	}
}
