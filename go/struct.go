package main

type Instance_of_instance struct {
	I    []Instance_Bdd
	Name string
}
type Instance_Bdd struct {
	I []all_bd
}

type all_bd struct {
	Smtp struct {
		Email string
		Pass  string
	}
	Categorie struct {
		Id   int
		Name string
	}

	Post struct {
		Id          int
		Id_cat      string
		Title_post  string
		Content     string
		Likes       string
		Posted_user string
		Last_Posted string
		Nb_Reply    string
	}

	User struct {
		Id              int
		Name            string
		Pswd            string
		Desc            string
		Email           string
		Profile_Picture string
		Rank_id         int
	}

	Temp_user struct {
		Id         int
		Name       string
		Email      string
		Pswd       string
		validation string
	}

	Profil struct {
		Id                  int
		User                string
		Joined              string
		Last_time_connected string
		Subjet_submit       string
		Email               string
		Desc                string
		Rank_id_profil      string
	}

	Comment struct {
		Id               int
		Id_post          string
		Date_comment     string
		User_posted      string
		Rank_User_Posted string
		Title_comment    string
		Reply_user       string
		Reply_user_rank  string
		Reply_content    string
		Content_comment  string
		Likes            string
	}
}
