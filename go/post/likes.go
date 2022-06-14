package post

import (
	"fmt"
	Database "forum/database"
	"net/http"
)

func add_like(w http.ResponseWriter, r *http.Request) {
	//check cookie (user connected)
	if r.FormValue("like") == "ok" {
		//check if post exist and is comment or post
		fmt.Print("add like")
		fmt.Print(r.FormValue(""))
		if Database.Select_column_where("post", "Id", r.FormValue("")) { //is a post
			//if post already not liked by the user

			//get nb of likes of post
			//else add like and add to table keep track of likes
			Database.Add_Like_To_DB("post", "Likes", "10", "Id", r.FormValue(""))
		}
		// } else if Database.Select_column_where("comment", "Id", r.FormValue("id_post")) { //is a comment
		// 	fmt.Print("is a comment")
		// }

	}
}
