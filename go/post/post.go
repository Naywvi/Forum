package post

type Post struct {
	Id                  int
	Id_cat              string
	Title               string
	Content             string
	Posted_user         string
	Posted_user_rank    string
	Likes               string
	Last_Posted         string
	Nb_Reply            string
	User_connected      string
	User_connected_rank string
}

type Comment struct {
	Id                  int
	Id_post             string
	Date_comment        string
	User_posted         string
	Rank_User_Posted    string
	Title_comment       string
	Reply_user          string
	Reply_user_rank     string
	Reply_content       string
	Content_comment     string
	Likes               string
	User_connected      string
	User_connected_rank string
}

type Statement_of_user struct {
	User      string
	Rank      string
	Post_Id   string
	Post      []Post
	Comment   []Comment
	Categorie string
}
