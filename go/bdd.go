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
	text := ""
	count := 0
	file, _ := os.Open(file_sql)
	scanner := bufio.NewScanner(file)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Func used ↓
func Checkalacon(result_table []string, input string) bool { //Check in lower case to be sure
	input = strings.ToLower(input)
	for i := range result_table {
		low := strings.ToLower(result_table[i])
		if low == input {
			return false
		}
	}
	return true
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

//#------------------------------------------------------------------------------------------------------------# ↓ Add to table func ↓

//Auto-Create Table
func Init_Database(database string, table_name string, txt string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Init db_test to start ↓
func categorie() {
	Inser_In_To_DB(Init_Database("dbtest.db", "categorie", Extract_File("../bdd/categorie_table.sql", 0, 2)), Extract_File("../bdd/categorie_table.sql", 8, 10), "categorie", Extract_File("../bdd/categorie_table.sql", 5, 6))
	fmt.Println("> Categorie Table was successfully created")
	fmt.Println("-> test_categorie was successfully created\n")
}

func post() {
	Inser_In_To_DB(Init_Database("dbtest.db", "post", Extract_File("../bdd/post_table.sql", 0, 6)), Extract_File("../bdd/post_table.sql", 13, 14), "post", Extract_File("../bdd/post_table.sql", 9, 10))
	fmt.Println("> Post Table was successfully created")
	fmt.Println("-> Test_post was successfully created\n")
}
func user() {
	Inser_In_To_DB(Init_Database("dbtest.db", "user", Extract_File("../bdd/user_table.sql", 0, 8)), Extract_File("../bdd/user_table.sql", 15, 16), "user", Extract_File("../bdd/user_table.sql", 11, 12))
	fmt.Println("> User Table was successfully created")
	fmt.Println("-> User test was successfully created | Login : Naywvi | pswd : 1230 |\n")
}
func rank() {
	for i := 1; i < 5; i++ {
		Inser_In_To_DB(Init_Database("dbtest.db", "rank", Extract_File("../bdd/rank_table.sql", 0, 10)), Extract_File("../bdd/rank_table.sql", 15+2*i, 16+2*i), "rank", Extract_File("../bdd/rank_table.sql", 13, 14))
	}
	fmt.Println("> Rank Table was successfully created\n")
}

func email_verification() {
	Inser_In_To_DB(Init_Database("dbtest.db", "email_owner", Extract_File("../bdd/email_verification_table.sql", 0, 2)), reveive_email_verification(), "email_owner", Extract_File("../bdd/email_verification_table.sql", 5, 6)) //<-- Os.Args email verification
}

//#------------------------------------------------------------------------------------------------------------# ↓ init bd ↓

//Init bdd
func InitBDD() {
	if _, err := os.Stat("./dbtest.db"); err == nil { //if bdd exist
		fmt.Println("The bdd is already here")

	} else if errors.Is(err, os.ErrNotExist) { //if bdd not exist > Re create
		email_verification()
		rank()
		user()
		categorie()
		post()
		fmt.Println("Bdd was successfully created, you are ready :)\n")
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ For register ↓
func Check_If_Exist(input string, check_field string, In_This_Table string) bool { //HORRIBLE !!!!
	var (
		db, _        = sql.Open("sqlite3", "dbtest.db")
		rows         = Select_Field_From_DB(db, check_field, In_This_Table)
		result_check []string
		u            User
	)

	for rows.Next() {
		err := rows.Scan(&u.Name)
		if err != nil {
			log.Fatal(err)
		}
		result_check = append(result_check, u.Name)
	}
	return Checkalacon(result_check, input)
}

func ADD_User_To_BDD(name string, pswd string, email string) {
	db, err := sql.Open("sqlite3", "dbtest.db") // init lg & ddb name
	if err != nil {
		log.Fatal(err)
	}
	var_send := []string{"'" + name + "',", "'" + pswd + "',", "'none_desc',", "'" + email + "',", "'none_picture',", "'3'"}
	Inser_In_To_DB(db, strings.Join(var_send, ""), "user", Extract_File("../bdd/user_table.sql", 11, 12))
}

//#------------------------------------------------------------------------------------------------------------# ↓ For login ↓
func Verif_If_Login_Exist(I *Instance, identifier string, pswd string, hash_pswd string) bool { //log
	for _, i := range I.I {
		if identifier == i.Email || identifier == i.Name {
			return CheckPasswordHash(pswd, hash_pswd)
		}
	}
	return false
}
func Check_If_Exist_Login(input_mail, input_pswd, hash_pswd string) bool { //Permet d'instancier User struct et de tout récup + check all cases
	var (
		I       = Instance{}
		u       = User{}
		db, _   = sql.Open("sqlite3", "dbtest.db")
		rows    = Select_All_From_DB(db, "user")
		input_m = strings.ToLower(input_mail)
		input_p = strings.ToLower(input_pswd)
	)

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Name, &u.Pswd, &u.Desc, &u.Email, &u.Profile_Picture, &u.Rank_id)
		u.Name = strings.ToLower(u.Name)
		u.Email = strings.ToLower(u.Email)
		I.I = append(I.I, u)
		if err != nil {
			log.Fatal(err)
		}
	}
	return Verif_If_Login_Exist(&I, input_m, input_p, hash_pswd)
}
