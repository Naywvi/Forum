package main

type Instance_of_instance struct {
	I []Instance_Bdd
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
		Id        int
		Id_catego string
		Name      string
		Contenu   string
		Likes     string
		User_id   string
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
}
