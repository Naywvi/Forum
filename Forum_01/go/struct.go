package main

type Categorie struct {
	Id   int
	Name string
}

type Post struct {
	Id        int
	Id_catego int
	Name      string
	Contenu   string
	Likes     string
	User_id   int
}

type Rank struct {
	Name      string
	Id        int
	Rank_perm struct {
		Rank_Id       int // /!\
		Name_rank     string
		Del_profile   int
		Signal_post   int
		Moove_post    int
		Comment_post  int
		Del_post      int
		Del_post_user int
		Create_post   int
		Admin         int
	}
}

type User struct {
	Id              int
	Name            string
	Pswd            string
	Desc            string
	Email           string
	Profile_Picture string
	Rank_id         int
}
