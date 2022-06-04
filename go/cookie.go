package main

import (
	"net/http"
	"time"
)

//#------------------------------------------------------------------------------------------------------------# ↓ cookie ↓

func SetCookie(w http.ResponseWriter, r *http.Request) {

	// 	Name:    "test",                  //name de la personne+mail+rank_id
	// 	Value:   "sssssssssdazdazdazdsa", //j_inscription
	expiration := time.Now().Add(time.Second * 900) //15 minutes
	cookie := http.Cookie{Name: "username", Value: "rank", Expires: expiration}
	http.SetCookie(w, &cookie)

}
