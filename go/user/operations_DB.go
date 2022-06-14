package user

import (
	"database/sql"
	Config "forum/config"
	"log"
	"strconv"
	"strings"

	Database "forum/database"
)

//Del user from table
func Del_User_From_Table(db *sql.DB, rows *sql.Rows, table, name_deleted, who_want string) { //all time send i of deleter
	var (
		Rows  []string
		u     = Config.All_bd{}
		marge = 0
		id    = ""
		index = 0
	)

	for rows.Next() {
		marge = 4 // <-- De combien je recule pour avoir l'id dans la table afin de le delect (vérification par la "Validation field")
		if who_want == "Validation" {
			err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.Validation)
			if err != nil {
				log.Fatal(err)
			}
			Rows = append(Rows, strconv.Itoa(u.Temp_user.Id), u.Temp_user.Name, u.Temp_user.Email, u.Temp_user.Pswd, u.Temp_user.Validation)
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
	if who_want == "Validation" {
		ADD_User_To_BDD(Rows[index-3], Rows[index-1], Rows[index-2], "3")
	}
}

//Add user to temp_user
func ADD_User_To_Temp(name, pswd, email, user_hash string) {
	var (
		db, err          = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		Default_user_arr = []string{"'" + name + "','" + pswd + "','" + email + "','" + user_hash + "'"}
		Default_user     = strings.Join(Default_user_arr, "")
	)
	if err != nil {
		log.Fatal(err)
	}
	Database.Inser_In_To_DB(db, Default_user, "temp_user", Database.Extract_File("../bdd/temp_user_table.sql", 8, 9))
}

//Add default user to bdd
func ADD_User_To_BDD(name, pswd, email, rank_id string) {
	var (
		db, err          = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		Default_user_arr = []string{"'" + name + "','" + pswd + "','none_desc','" + email + "','none_picture','" + rank_id + "'"}
		Default_user     = strings.Join(Default_user_arr, "")
	)

	if err != nil {
		log.Fatal(err)
	}

	Database.Inser_In_To_DB(db, Default_user, "user", Database.Extract_File("../bdd/user_table.sql", 11, 12))
	New_Profil(name, email, rank_id)
}
