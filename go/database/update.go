package database

import (
	"database/sql"
	"fmt"
	Config "forum/config"
	"log"
	"strconv"
)

func Update_Field(Table, set_attribute, set_attribute_value, where_attribute, where_attribute_value string) {
	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("UPDATE " + Table + " SET " + set_attribute + " = '" + set_attribute_value + "' WHERE " + where_attribute + " = '" + where_attribute_value + "';")
}

func Add_Like_To_DB(Table, set_attribute, set_attribute_value, where_attribute, where_attribute_value string) {
	var (
		db, err = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		like, _ = strconv.Atoi(set_attribute_value)
	)
	if err != nil {
		log.Fatal(err)
	}
	like += 1
	fmt.Print("UPDATE " + Table + " SET " + set_attribute + " = " + strconv.Itoa(like) + " WHERE " + where_attribute + " = " + where_attribute_value + ";")
	db.Exec("UPDATE " + Table + " SET " + set_attribute + " = " + strconv.Itoa(like) + " WHERE " + where_attribute + " = " + where_attribute_value + " ;")
}
