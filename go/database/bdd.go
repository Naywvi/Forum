package database

import (
	"bufio"
	"database/sql"
	"log"
	"os"

	Config "forum/config"
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

func Print_Rows(rows *sql.Rows, table string) []Config.All_bd {
	var (
		I Config.Instance_Bdd
		u = Config.All_bd{}
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
			err := rows.Scan(&u.Temp_user.Id, &u.Temp_user.Name, &u.Temp_user.Email, &u.Temp_user.Pswd, &u.Temp_user.Validation)
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

//#------------------------------------------------------------------------------------------------------------# ↓ Init Add to table ↓

//Auto-Create Table
func Init_Database(table_name, txt string) *sql.DB {
	db, err := sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
	if err != nil {
		log.Fatal(err)
	}
	sql := "CREATE TABLE IF NOT EXISTS " + table_name + "(" + txt + ")"
	db.Exec(sql)
	return db
}
