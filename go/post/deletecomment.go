package post

import (
	"database/sql"
	Config "forum/config"
	"net/http"
	"strconv"
)

func delete_comment(w http.ResponseWriter, r *http.Request, deletePostid string) {
	db, _ := sql.Open(Config.Bdd.Langage, Config.Bdd.Name)

	db.Exec("DELETE FROM post WHERE id = " + deletePostid)                // delete post
	db.Exec("DELETE FROM comment WHERE id_post = '" + deletePostid + "'") // delete all comment of this post

	//get nb of reply of post
	var nb_reply int
	db.QueryRow("SELECT Nb_Reply FROM post WHERE id = " + deletePostid).Scan(&nb_reply) // get the number of reply of the post
	nb_reply--                                                                          // decrement nb of reply

	db.Exec("UPDATE FROM post WHERE id_post = '" + strconv.Itoa(nb_reply) + "'") // update the number of reply of the post
}
