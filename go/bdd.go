package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

//Extract sql-file & return it (select interval in file with end/ start)
func Extract_File(file_sql string, start int, end int) string {
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

//Add default user to bdd
func ADD_User_To_BDD(name, pswd, email string) {
	var (
		db, err          = sql.Open(Bdd.Langage, Bdd.Name)
		Default_user_arr = []string{"'" + name + "',", "'" + pswd + "',", "'none_desc',", "'" + email + "',", "'none_picture',", "'3'"}
		Default_user     = strings.Join(Default_user_arr, "")
	)

	if err != nil {
		log.Fatal(err)
	}

	Inser_In_To_DB(db, Default_user, "user", Extract_File("../bdd/user_table.sql", 11, 12))
}

//#------------------------------------------------------------------------------------------------------------# ↓ Select on table ↓

func Select_All_From_DB(db *sql.DB, table string) *sql.Rows {
	result, _ := db.Query("SELECT * FROM " + table)
	return result
}
func Select_Field_From_DB(db *sql.DB, field string, table string) *sql.Rows {
	result, _ := db.Query("SELECT " + field + " FROM " + table)
	return result
}

//#------------------------------------------------------------------------------------------------------------# ↓ Init Add to table ↓

//Auto-Create Table
func Init_Database(table_name string, txt string) *sql.DB {
	db, err := sql.Open(Bdd.Langage, Bdd.Name)
	if err != nil {
		log.Fatal(err)
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table_name + "(" + txt + ")"
	db.Exec(sql)
	return db
}

//Auto increment field & Value on table
func Inser_In_To_DB(db *sql.DB, var_receive string, table_name string, table_field string) (int64, error) { //
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

		fmt.Println("---------------------")

		return "'" + username + "','" + pswd + "','none_desc','" + mail + "','none_picture','3'"

	} else if who_want == "email_verification_table" {

		fmt.Println("Write your mail to send request register.")
		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password -> ")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n') //<-- Email need \n to connect
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		fmt.Println("---------------------")

		return "'" + mail + "','" + pswd + "'"
	} else if who_want == "Bdd_Name" {
		fmt.Println("Choose name of your Data_Base.")
		Name_Bd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		Bdd.Name = strings.TrimSpace(Name_Bd) + ".db"
		Bdd.Langage = "sqlite3"
	}
	return ""
}

//Simple print on shell
func Is_Ok(Printable string, Second_Printable string) {

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
	Inser_In_To_DB(Init_Database("post", Extract_File("../bdd/post_table.sql", 0, 6)), Extract_File("../bdd/post_table.sql", 13, 14), "post", Extract_File("../bdd/post_table.sql", 9, 10))
	Is_Ok("Post", "Test_post")
}
func user() {
	Inser_In_To_DB(Init_Database("user", Extract_File("../bdd/user_table.sql", 0, 8)), Terminal_Init_Table("add_user_table"), "user", Extract_File("../bdd/user_table.sql", 11, 12))
	Is_Ok("User", "")
}
func rank() {
	for i := 1; i < 5; i++ {
		Inser_In_To_DB(Init_Database("rank", Extract_File("../bdd/rank_table.sql", 0, 10)), Extract_File("../bdd/rank_table.sql", 15+2*i, 16+2*i), "rank", Extract_File("../bdd/rank_table.sql", 13, 14))
	}
	Is_Ok("Rank", "")
}

func email_verification() {
	Inser_In_To_DB(Init_Database("email_owner", Extract_File("../bdd/email_verification_table.sql", 0, 2)), Terminal_Init_Table("email_verification_table"), "email_owner", Extract_File("../bdd/email_verification_table.sql", 5, 6)) //<-- Os.Args email verification
}

//#------------------------------------------------------------------------------------------------------------# ↓ init bd ↓

//Init bdd
func InitBDD() {
	if _, err := os.Stat("./" + Bdd.Name); err == nil { //<-- If bdd exist
		fmt.Println("The bdd, " + Bdd.Name + " is already here")

	} else if errors.Is(err, os.ErrNotExist) { //<-- If bdd not exist > Re create
		fmt.Println("--->Create a new db<--- ?")
		email_verification()
		rank()
		user()
		categorie()
		post()
		fmt.Println("Bdd, " + Bdd.Name + " was successfully created, you are ready :)\n")
	}
}
