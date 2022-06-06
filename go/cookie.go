package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Encrypt / Decrypt cookie ↓

//Encrypte Cookie
func Encrypt_Cookie(input string) string { //Cookie n'aime pas certains caract base 64 donc filtre "=, ,$,[]"
	var (
		Cookie = base64.StdEncoding.EncodeToString([]byte(input))
		str    = ""
	)

	for i := 0; i < len(Cookie); i++ {
		if (Cookie[i] >= 'A' && Cookie[i] <= 'Z') || (Cookie[i] >= 'a' && Cookie[i] <= 'z') {
			str += string(Cookie[i])
		}
	}

	return string(str)

}

//Decrypte Cookie
func Decrypt_Cookie(input, input_hash string) bool {

	result := Encrypt_Cookie(input)
	if result == input_hash {
		return true
	} else {
		return false
	}
}

//#------------------------------------------------------------------------------------------------------------# ↓ Check if Cookie is Correct ↓

//Check User cookie + return rank value
func Check_Cookie(w http.ResponseWriter, r *http.Request) (bool, string) {
	var (
		Cookie   = r.Cookies()
		str_id   = ""
		str_hash = ""
		count    = 0
		result1  = false
		result2  = false
	)

	for _, c := range Cookie {

		for _, i := range c.Name {

			if i == '_' {
				count = 1
			} else if count == 0 {
				str_id += string(i)
			} else if count == 1 {
				str_hash += string(i)
			}
		}

		result1 = Decrypt_Cookie(str_id, str_hash) //<-- Name

		var (
			str_id   = ""
			str_hash = ""
			count    = 0
		)

		for _, i := range c.Value {

			if i == '_' {
				count = 1
			} else if count == 0 {
				str_id += string(i)
			} else if count == 1 {
				str_hash += string(i)
			}
		}

		result2 = Decrypt_Cookie(str_id, str_hash) //<-- value(-> rank_id)

		if result1 == true && result2 == true {
			add(w, r) //<-- Init + 15 min
			return true, str_id
		} else {
			del(w, r) //<-- delete
			return false, "4"
		}

	}
	return false, "4"
}

//#------------------------------------------------------------------------------------------------------------# ↓ Manage cookie ↓

//add +15 min actual cookie
func add(w http.ResponseWriter, r *http.Request) {
	Cookie := r.Cookies()
	for _, c := range Cookie {
		c.Expires = time.Now().Add(time.Second * 900) //+15 min
		http.SetCookie(w, c)
	}
}

//delete actual
func del(w http.ResponseWriter, r *http.Request) {
	Cookie := r.Cookies()
	for _, c := range Cookie {
		c.Expires = time.Now()
		http.SetCookie(w, c)
	}

}

//#------------------------------------------------------------------------------------------------------------# ↓ Set cookie ↓

//Set the Cookie
func SettCookie(w http.ResponseWriter, r *http.Request) {

	expiration := time.Now().Add(time.Second * 900) //15 minutes

	Name := Connected.User + "_" + Connected.User_Hased //<-- Set personnal token
	Value := Connected.Rank_Id + "_" + Connected.Rank_Id_Hashed
	cookie := http.Cookie{Name: Name, Value: Value, Expires: expiration}

	http.SetCookie(w, &cookie)

}
