//forumnlt@gmail.com	12301230789Aa
package main

import (
	"bufio"
	"fmt"
	"os"
)

//#------------------------------------------------------------------------------------------------------------# ↓ init mail and hash password ↓
// Call in BDD.go & send hash password
func reveive_email_verification() string {
	fmt.Println("---------------------")
	fmt.Print("\n/!\\End email with > '\nemail -> ")
	mail, _ := bufio.NewReader(os.Stdin).ReadString('\'')
	fmt.Print("\n/!\\End mdp with > '\nmdp -> ")
	pswd, _ := bufio.NewReader(os.Stdin).ReadString('\'')
	fmt.Println("---------------------")
	sz := len(pswd)

	if sz > 0 && pswd[sz-1] == '\'' {
		pswd = pswd[:sz-1]
	}
	pswd, _ = HashPassword(pswd)
	return "'" + mail + ",'" + pswd + "'"
}
