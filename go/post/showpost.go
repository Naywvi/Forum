package post

import (
	"database/sql"
	"fmt"
	Config "forum/config"
	Database "forum/database"
	User "forum/user"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func Show_Post(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("")

	add_like(w, r)

	var (
		_, statement, User = User.Check_Cookie(w, r)
		db, err            = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
		time               = time.Now()
		time_str           = time.String()
		pos                = Statement_of_user{}
		rows               = Database.Select_column("post", "Id", query)
		rows_comment       = Database.Select_column("comment", "Id_post", query)
		instance           Config.All_bd
		POST               = Post{}
		COMMENT            = Comment{}
	)
	pos.User = User
	pos.Rank = statement
	pos.Post_Id = query
	//#----------------------------------------------------------------------------#
	//<<< add post
	for rows.Next() {
		err := rows.Scan(&instance.Post.Id, &instance.Post.Id_cat, &instance.Post.Title_post, &instance.Post.Content, &instance.Post.Likes, &instance.Post.Posted_user, &instance.Post.Last_Posted, &instance.Post.Nb_Reply)

		if err != nil {
			log.Fatal(err)
		}

		POST.Id = instance.Post.Id
		POST.Id_cat = instance.Post.Id_cat
		POST.Title = instance.Post.Title_post
		POST.Content = instance.Post.Content
		POST.Posted_user = instance.Post.Posted_user

		result := ""
		req := Database.Select_column("user", "name", POST.Posted_user)
		for req.Next() {
			errf := req.Scan(&instance.User.Id, &instance.User.Name, &instance.User.Pswd, &instance.User.Desc, &instance.User.Email, &instance.User.Profile_Picture, &instance.User.Rank_id)
			if errf != nil {
				log.Fatal(errf)
			}
			result = strconv.Itoa(instance.User.Rank_id)

		}
		POST.Posted_user_rank = result

		POST.Likes = instance.Post.Likes
		POST.Last_Posted = instance.Post.Last_Posted
		POST.Nb_Reply = instance.Post.Nb_Reply
		POST.User_connected = User
		POST.User_connected_rank = statement

		//<<Append the post
		pos.Post = append(pos.Post, POST)

	}
	//#----------------------------------------------------------------------------#
	//<<add all comment by post id

	for rows_comment.Next() {

		rows_comment.Scan(&instance.Comment.Id, &instance.Comment.Id_post, &instance.Comment.Date_comment, &instance.Comment.User_posted, &instance.Comment.Rank_User_Posted, &instance.Comment.Title_comment, &instance.Comment.Reply_user, &instance.Comment.Reply_user_rank, &instance.Comment.Reply_content, &instance.Comment.Content_comment, &instance.Comment.Likes)
		//<< add post
		COMMENT.Id = instance.Comment.Id
		COMMENT.Id_post = instance.Comment.Id_post
		COMMENT.Date_comment = instance.Comment.Date_comment
		COMMENT.User_posted = instance.Comment.User_posted
		COMMENT.Rank_User_Posted = instance.Comment.Rank_User_Posted
		COMMENT.Title_comment = instance.Comment.Title_comment
		COMMENT.Reply_user = instance.Comment.Reply_user
		COMMENT.Reply_user_rank = instance.Comment.Reply_user_rank
		COMMENT.Reply_content = instance.Comment.Reply_content
		COMMENT.Content_comment = instance.Comment.Content_comment
		COMMENT.Likes = instance.Comment.Likes
		COMMENT.User_connected = User
		COMMENT.User_connected_rank = statement
		//<< add post
		//<<Append the post
		pos.Comment = append(pos.Comment, COMMENT)

	}

	//<<
	if r.Method == "GET" {
		var (
			deletePostid = r.FormValue("")
			deletePost   = r.FormValue("deletep")
			deletec      = r.FormValue("deletec")
		)

		//<<< add post

		if deletePost == deletePostid { //Delete post User-poster or admin
			if (len(deletePost) > 0 && statement == "1") || (len(deletePost) > 0 && statement == "2") || (len(deletePost) > 0 && (User == POST.Posted_user)) {
				db.Exec("DELETE FROM post WHERE id = " + deletePostid)
				db.Exec("DELETE FROM comment WHERE id_post = '" + deletePostid + "'")
			}
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
		} else if (len(deletec) > 0 && deletec == COMMENT.User_posted) || (len(deletec) > 0 && (statement == "1" || statement == "2")) {
			db.Exec("DELETE FROM comment WHERE Id = '" + deletec + "'")
			fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/"; </script>`)
		}

		template.Must(template.ParseFiles(filepath.Join(Config.TemplatesDir, "../static/templates/post.html"))).Execute(w, pos)

	} else if r.Method == "POST" {
		r.ParseForm()
		var (
			Comment_content_parse     = r.Form["Comment_Content"][0]
			Comment_content_parse_sql = strings.Replace(Comment_content_parse, "'", "’", 10000) //<< Replace ' to > `  protect from sql_exploit
			reply_to                  = r.FormValue("Reply_to")
		)
		if err != nil {
			log.Fatal(err)
		}
		//add reply
		t, _ := strconv.Atoi(POST.Nb_Reply)
		t += 1
		POST.Nb_Reply = strconv.Itoa(t)
		db.Exec("UPDATE post SET Nb_Reply = '" + POST.Nb_Reply + "' WHERE Id = " + strconv.Itoa(POST.Id) + ";")
		//

		if reply_to == "post" { // reply to post >> push comment

			var (
				var_p    = []string{"'" + pos.Post_Id + "','" + time_str[0:10] + "','" + User + "','" + statement + "','nil','no_reply','no_reply','no_reply','" + Comment_content_parse_sql + "','0'"}
				var_pstr = strings.Join(var_p, "")
				db, _    = sql.Open(Config.Bdd.Langage, Config.Bdd.Name)
			)

			Database.Inser_In_To_DB(db, var_pstr, "comment", Database.Extract_File("../bdd/comment_table.sql", 14, 15)) //<-- Push the post with no reply

		} else if len(reply_to) > 0 { //reply to user on the post >> push comment
			var (
				// Check_user_exist = Check_If_Exist(reply_to, "", "Name", "user", "Register")

				reply_content = r.Form["Reply_content_"+reply_to][0]

				reply_user      = r.Form["Reply_USER_"+reply_to][0]
				reply_user_rank = r.Form["Reply_USER_RANK_"+reply_to][0]
				var_p           = []string{"'" + pos.Post_Id + "','" + time_str[0:10] + "','" + User + "','" + statement + "','nil','" + reply_user + "','" + reply_user_rank + "','" + reply_content + "','" + Comment_content_parse_sql + "','0'"}
				var_pstr        = strings.Join(var_p, "")
			)

			// if Check_user_exist == false {
			//fmt.Fprint(w, "<script>alert('Wrong user who to reply to'); </script>")

			Database.Inser_In_To_DB(db, var_pstr, "comment", Database.Extract_File("../bdd/comment_table.sql", 14, 15)) //<-- Push the post with reply

		} else {

			Config.Send_Error(w, r)
			return
		}
		fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/post?=`+pos.Post_Id+`&Reply_to=post"; </script>`)

	} else {
		Config.Send_Error(w, r)

		return
	}
}
