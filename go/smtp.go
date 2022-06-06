//forumnlt@gmail.com	12301230789Aa
package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"text/template"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

//email / password (non hash) / >inscription(récup email) />reset password(email > reset pass) />alert (>all bdd)
func Register_Smtp(email_register, Name_User string) {
	var (
		subject_mail = "Forum NLT - Nouveau membre !"
		to           = []string{email_register}
		Who_Want     = "Register"
	)

	Init_Smtp(to, subject_mail, Name_User, Who_Want)
}

func Reset_Pswd_Smtp(email_reset string) {

}

// func alert_Smtp() {
// 	smtp() //<-- Envoie tout les emails
// }
func Init_Smtp(to []string, subject_mail, Name_User, Who_Want string) { //<-- Récup log bdd (retirer hash)

	// Sender data.
	from := "forum.nlt@hotmail.com" //à recup
	password := "12301230789Aa"

	// Receiver email address.

	// smtp server configuration.
	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	conn, err := net.Dial("tcp", "smtp.office365.com:587")
	if err != nil {
		println(err)
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println(err)
	}

	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
	}

	auth := LoginAuth(from, password)

	if err = c.Auth(auth); err != nil {
		println(err)
	}

	if Who_Want == "Register" {
		t, _ := template.ParseFiles("../static/templates/smtp/register.html")
		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf(""+subject_mail+" \n%s\n\n", mimeHeaders)))
		//<<<<<
		t.Execute(&body, struct {
			Name string
		}{
			Name: Name_User, //<--- Name de la personne
		})
		//<<<<<
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Sending email.

	fmt.Println("Email Sent!")
}
