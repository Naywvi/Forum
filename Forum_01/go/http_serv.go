package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var templatesDir = os.Getenv("TEMPLATES_DIR")

func HomePage(w http.ResponseWriter, r *http.Request) {

	template.Must(template.ParseFiles(filepath.Join(templatesDir, "../static/index.html"))).Execute(w, " ")
}

func EncryptPassword(password string) (string, error) { // fonction qui nous permettra d'encrypter nos mots de passe avant de les intégrer à la db
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool { // fonction qui nous permet de vérifier si le mot de passe rentré correspond au mdp encrypté stocké dans la database
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) // ici on encrypt l'entrée utilisateur et on la compare au mdp stocké dans la db, encrypté lui aussi
	return err == nil
}

func httpServ() {

	fs := http.FileServer(http.Dir("../static")) // <- ce qu'on envoie en static vers le serv
	http.Handle("/", fs)

	http.ListenAndServe(":1010", nil)
	fmt.Print("Started https serv successfully on http://localhost:1010")
}
