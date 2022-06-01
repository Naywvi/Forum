package main

import (
	"database/sql"
	"log"
)

func initDatabase(database string) *sql.DB { // ici on créé une fonction qui initialise une nouvelle database
	db, err := sql.Open("sqlite3", database) // on ouvre une nouvelle database avec sqlite3 qu'on nomme par le paramètre qu'on donne à notre fonction
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
				PRAGMA foreign_keys = ON;
				
				CREATE TABLE IF NOT EXISTS rank (
					Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					Name TEXT NOT NULL
					Create_post INTEGER NOT NULL
					Del_profile INTEGER NOT NULL
					Signal_post INTEGER NOT NULL
					Moove_post INTEGER NOT NULL
					Comment_post INTEGER NOT NULL
					Del_post INTEGER NOT NULL
					Del_post_user INTEGER NOT NULL
					Admin INTEGER NOT NULL
				);

				CREATE TABLE IF NOT EXISTS user (
					Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					Name TEXT NOT NULL
					Pswd TEXT NOT NULL
					Desc TEXT NOT NULL
					Email TEXT NOT NULL
					Profile_Picture TEXT NOT NULL
					Rank_id INTEGER NOT NULL,
					FOREIGN KEY (Rank_id) REFERENCES rank(Id)
				);

				CREATE TABLE IF NOT EXISTS categorie (
					Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					Name TEXT NOT NULL
				);

				CREATE TABLE IF NOT EXISTS post (
					Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					Id_catego INTEGER NOT NULL,
					FOREIGN KEY (Id_catego) REFERENCES categorie(Id),
					Name TEXT NOT NULL
					Contenu TEXT NOT NULL
					Likes INTEGER NOT NULL
					User_id INTEGER NOT NULL,
					FOREIGN KEY (User_id) REFERENCES user(Id)
				);
				`
	//PRAGMA foreign_keys = ON; permet de donner à sql l'accès aux tables reliers entres elles
	//FOREIGN KEY permet de relierles tables entre elles (type_id) pointe vers > (id) de la table types
	db.Exec(sqlStmt) // On affiche l'erreur si ça merde
	return db
}

func insertIntoRank(db *sql.DB, name string, create_post int, del_profile int, signal_post int, moove_post int, comment_post int, del_post_user int, admin int) (int64, error) { // Ici une fonction qui nous permet d'insérer un nouveau nom  de rang dans la table "rank" de notre database
	result, _ := db.Exec(`INSERT INTO rank (Name,Create_post,Del_profile,Signal_post,Moove_post,Comment_post,Del_post_user,Admin) VALUES (?,?,?,?,?,?,?,?)`, name, create_post, del_profile, signal_post, moove_post, comment_post, del_post_user, admin)
	return result.LastInsertId()
}

func insertIntoUser(db *sql.DB, name string, pswd string, desc string, email string, profile_picture string, rank_id int) (int64, error) { // Ici une fonction qui nous permet d'insérer un nouveau type dans la table "Types"(ou autres) de notre database
	result, _ := db.Exec(`INSERT INTO user (Name,Pswd,Desc,Email,Profile_picture,Rank_id) VALUES (?,?,?,?,?,?)`, name, pswd, desc, email, profile_picture, rank_id)
	return result.LastInsertId()
}

func insertIntoPost(db *sql.DB, id_catego int, name string, contenu string, likes int, user_id int) (int64, error) { // Ici une fonction qui nous permet d'insérer un nouveau type dans la table "Types"(ou autres) de notre database
	result, _ := db.Exec(`INSERT INTO user (Id_catego,Name,Contenu,Likes,User_id) VALUES (?,?,?,?,?)`, id_catego, name, contenu, likes, user_id)
	return result.LastInsertId()
}

func insertIntoCategorie(db *sql.DB, name string) (int64, error) { // Ici une fonction qui nous permet d'insérer un nouveau type dans la table "Types"(ou autres) de notre database
	result, _ := db.Exec(`INSERT INTO user (Name) VALUES (?)`, name)
	return result.LastInsertId()
}
