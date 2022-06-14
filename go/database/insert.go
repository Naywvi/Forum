package database

import (
	"database/sql"
	"log"
)

//Auto increment field & Value on table
func Inser_In_To_DB(db *sql.DB, var_receive, table_name, table_field string) (int64, error) { //
	result, err := db.Exec("INSERT INTO " + table_name + " (" + table_field + ")" + " VALUES (" + var_receive + ")")
	if err != nil {
		log.Fatal(err)
	}
	return result.LastInsertId()
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
	Inser_In_To_DB(Init_Database("comment", Extract_File("../bdd/comment_table.sql", 0, 11)), Extract_File("../bdd/comment_table.sql", 18, 19), "comment", Extract_File("../bdd/comment_table.sql", 14, 15))
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
	Inser_In_To_DB(Init_Database("profil", Extract_File("../bdd/profil_table.sql", 0, 8)), Extract_File("../bdd/profil_table.sql", 15, 16), "profil", Extract_File("../bdd/profil_table.sql", 11, 12))
	Is_Ok("profil", "")
}
