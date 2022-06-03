package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

//Extract sql-file & return it (select interval in file with end/ start)
func SqlExtract(file_sql string, start int, end int) string {
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

//#------------------------------------------------------------------------------------------------------------#

//Auto-Create Table
func initDatabase(database string, table_name string, txt string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table_name + "(" + txt + ")"
	db.Exec(sql)
	return db
}

//Auto increment field & Value on table
func InsterInTo(db *sql.DB, var_receive string, table_name string, table_field string) (int64, error) { //
	result, err := db.Exec("INSERT INTO " + table_name + " (" + table_field + ")" + " VALUES (" + var_receive + ")")
	if err != nil {
		log.Fatal(err)
	}
	return result.LastInsertId()
}

//#------------------------------------------------------------------------------------------------------------#
func categorie() {
	InsterInTo(initDatabase("dbtest.db", "categorie", SqlExtract("../bdd/categorie_table.sql", 0, 2)), SqlExtract("../bdd/categorie_table.sql", 8, 10), "categorie", SqlExtract("../bdd/categorie_table.sql", 5, 6))
	fmt.Println("> Categorie Table was successfully created")
	fmt.Println("-> test_categorie was successfully created\n")
}
func post() {
	InsterInTo(initDatabase("dbtest.db", "post", SqlExtract("../bdd/post_table.sql", 0, 6)), SqlExtract("../bdd/post_table.sql", 13, 14), "post", SqlExtract("../bdd/post_table.sql", 9, 10))
	fmt.Println("> Post Table was successfully created")
	fmt.Println("-> Test_post was successfully created\n")
}
func user() {
	InsterInTo(initDatabase("dbtest.db", "user", SqlExtract("../bdd/user_table.sql", 0, 8)), SqlExtract("../bdd/user_table.sql", 15, 16), "user", SqlExtract("../bdd/user_table.sql", 11, 12))
	fmt.Println("> User Table was successfully created")
	fmt.Println("-> User test was successfully created | Login : Naywvi | pswd : 1230 |\n")
}
func rank() {
	for i := 1; i < 5; i++ {
		InsterInTo(initDatabase("dbtest.db", "rank", SqlExtract("../bdd/rank_table.sql", 0, 10)), SqlExtract("../bdd/rank_table.sql", 15+2*i, 16+2*i), "rank", SqlExtract("../bdd/rank_table.sql", 13, 14))
	}
	fmt.Println("> Rank Table was successfully created\n")
}

/*
Exemple create / insert : bdd_name | table |field
SqlExtract > Récup dans le fichier avec un interval
> Renvoie à InserInTo([Nom base de donnée].db , '[Nom de la table voulu]')
> Renvoie à initDataBase([db] , [Var recup de sqlextract()] , [name table] , [var field récup de sqlextract()])
tout en une ligne ;)
*/

//#------------------------------------------------------------------------------------------------------------#

//Init bdd
func InitBDD() {
	if _, err := os.Stat("./dbtest.db"); err == nil { //if bdd exist
		fmt.Println("The bdd is already here")

	} else if errors.Is(err, os.ErrNotExist) { //if bdd not exist > Re create
		rank()
		user()
		categorie()
		post()
		fmt.Println("Bdd was successfully created, you are ready :)\n")
	}

}
