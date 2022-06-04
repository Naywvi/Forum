//forumnlt@gmail.com	12301230789Aa
package main

//strings.TrimSpace(VARIABLE) <-- remove tabulation

//#------------------------------------------------------------------------------------------------------------# ↓ init mail and hash password ↓
// Call in BDD.go & send hash password

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
