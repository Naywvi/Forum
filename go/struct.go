package main

type Instance_Bdd struct {
	I []all_bd
}

type all_bd struct {
	Categorie struct {
		Id   int
		Name string
	}

	Post struct {
		Id        int
		Id_catego int
		Name      string
		Contenu   string
		Likes     string
		User_id   int
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
}
