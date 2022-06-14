package database

import (
	"database/sql"
	"fmt"
	Config "forum/config"
	"log"
)

func Select_column_where(Table_name, Table_field, input string) bool { // select columns where field = input and returns true or false
	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Query("SELECT * FROM " + Table_name + " WHERE " + Table_field + " = " + input + ";")
	fmt.Print(Config.HandleError(err))
	return result != nil
}

func Select_column(Table_name, Table_field, input string) *sql.Rows { //only string
	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Query("SELECT * FROM " + Table_name + " WHERE " + Table_field + " = '" + input + "';")
	fmt.Print(Config.HandleError(err))
	return result
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
func Select_All_Rows_Table(db *sql.DB, table []string) Config.Instance_of_instance {
	var (
		I   Config.Instance_Bdd
		I_I Config.Instance_of_instance
	)
	for i := range table {
		I.I = Print_Rows(Select_All_From_DB(db, table[i]), table[i])
		I_I.I = append(I_I.I, I) // <-- To send on one template
	}
	return I_I
}
