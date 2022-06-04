//forumnlt@gmail.com	12301230789Aa
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//strings.TrimSpace(VARIABLE) <-- remove tabulation

//#------------------------------------------------------------------------------------------------------------# ↓ init mail and hash password ↓
// Call in BDD.go & send hash password
func Terminal_Init_Table(who_whant string) string {

	fmt.Println("---------------------")

	if who_whant == "add_user_table" {

		fmt.Println("Create a default user to connect on the web site.")
		fmt.Print("username ->")

		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password ->")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		fmt.Println("---------------------")

		return "'" + username + "','" + pswd + "','none_desc','" + mail + "','none_picture','3'"

	} else if who_whant == "email_verification_table" {

		fmt.Println("Write your mail to send request register.")
		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password -> ")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n') //<-- Email need \n to connect
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		fmt.Println("---------------------")

		return "'" + mail + "','" + pswd + "'"
	}

	return ""
}

//#------------------------------------------------------------------------------------------------------------# ↓ Send mail Register ↓
// func SendMailRegister() {

// 	from := "naj_lak93@hotmail.fr"
// 	password := "123012301230789Aa"

// 	toEmailAddress := "naj-lak93@hotmail.fr"
// 	to := []string{toEmailAddress}

// 	host := "smtp.gmail.com"
// 	port := "587"
// 	address := host + ":" + port

// 	subject := "Subject: This is the subject of the mail\n"
// 	body := "This is the body of the mail"
// 	message := []byte(subject + body)

// 	auth := smtp.PlainAuth("", from, password, host)

// 	err := smtp.SendMail(address, auth, from, to, message)
// 	if err != nil {
// 		panic(err)
// 	}
// }
