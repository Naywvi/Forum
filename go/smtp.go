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

//#------------------------------------------------------------------------------------------------------------# ↓ Register -> verification mail smtp ↓

//email / password (non hash) / >inscription(récup email) />reset password(email > reset pass) />alert (>all bdd)
func Register_Smtp(email_register, Name_User, user_hash string) {
	var (
		to       = []string{email_register}
		Who_Want = "Register"
	)

	Init_Smtp(to, Name_User, user_hash, Who_Want)
}

//#------------------------------------------------------------------------------------------------------------# ↓ Reset pass -> send url smtp ↓
func Reset_Pswd_Smtp(email_reset string) {

}

//#------------------------------------------------------------------------------------------------------------# ↓ Send alert to all users smtp ↓

func alert_Smtp(to []string, path string) {
	Init_Smtp(to, path, "", "alert")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Manage smtp ↓
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
func Init_Smtp(to []string, Name_User, user_hash, Who_Want string) { //<-- Récup log bdd (retirer hash)
	//forum.nlt@hotmail.com
	// Sender data.
	from := "naj_lak93@hotmail.fr" //à recup
	password := "123012301230789Aa"

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
		t, _ := template.ParseFiles("../static/templates/smtp/confirm_verification_mail.html")
		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject : Forum NLT - Nouveau membre ! \n%s\n\n", mimeHeaders)))
		//<<<<<
		t.Execute(&body, struct {
			Name string
			Link string
		}{
			Name: Name_User, //<--- Name de la personne
			Link: user_hash,
		})
		//<<<<<
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if Who_Want == "alert" {
		alert_path := Name_User // flemme
		t, _ := template.ParseFiles(alert_path)
		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject : Forum NLT - Nlt-Bot (alert) ! \n%s\n\n", mimeHeaders)))
		//<<<<<
		t.Execute(&body, struct {
			Name string
		}{
			Name: "à toi", //<--- Name de la personne
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
