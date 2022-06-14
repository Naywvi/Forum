package database

import (
	"bufio"
	"errors"
	"fmt"
	Config "forum/config"
	"os"
	"strings"
)

//Init bdd
func InitBDD() {
	if _, err := os.Stat("./" + Config.Bdd.Name); err == nil { //<-- If bdd exist
		fmt.Println("The bdd, " + Config.Bdd.Name + " is already here")

	} else if errors.Is(err, os.ErrNotExist) { //<-- If bdd not exist > Re create
		fmt.Println("--->Create a new db<--- ?")
		email_verification()
		temp_user()
		user()
		categorie()
		post()
		comment()
		profilt()
		fmt.Println("Config.Bdd, " + Config.Bdd.Name + " was successfully created, you are ready :)\n")
	}
}

//Selection input of shell during init bdd
func Terminal_Init_Table(who_want string) string {

	fmt.Println("---------------------")

	if who_want == "add_user_table" {

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
		fmt.Println("User -> 3 | Moderator -> 2 | Admin -> 1")
		fmt.Print("rank_id ->")

		rank, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		rank = strings.TrimSpace(rank)
		fmt.Println("---------------------")

		return "'" + username + "','" + pswd + "','none_desc','" + mail + "','none_picture','" + rank + "'"

	} else if who_want == "email_verification_table" {

		fmt.Println("Write your mail to send request register.")
		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)

		fmt.Print("password -> ")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n') //<-- Email need \n to connect
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		return "'" + mail + "','" + pswd + "'"
	} else if who_want == "Bdd_Name" {
		fmt.Println("Choose name of your Data_Base.")
		Name_Bd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		Config.Bdd.Name = strings.TrimSpace(Name_Bd) + ".db"
		Config.Bdd.Langage = "sqlite3"

		fmt.Println(Config.Bdd.Name, "Is the database selected. Do you validate ? [y] [n] ")

		Check, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		Check = strings.TrimSpace(Check)

		if Check == "n" || Check == "N" || Check == "no" || Check == "NO" {
			return Terminal_Init_Table(who_want)
		} else if Check == "y" || Check == "Y" || Check == "YES" || Check == "yes" {
			fmt.Println("You are ready now")
		} else {
			fmt.Println("Wrong selection")
			fmt.Println("---------------------")
			return Terminal_Init_Table(who_want)
		}

	} else if who_want == "temp_user" {
		fmt.Println("Create a temp user to check the Validation table.")
		fmt.Print("username ->")

		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("email ->")

		mail, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		mail = strings.TrimSpace(mail)
		mail_hash := InitHashPswd(mail)

		fmt.Print("password ->")

		pswd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		pswd = strings.TrimSpace(pswd)
		pswd, _ = HashPassword(pswd)

		fmt.Println("---------------------")
		return "'" + username + "','" + pswd + "','" + mail + "','" + mail_hash + "'"
	}
	fmt.Println("---------------------")
	return ""
}

//Simple print on shell
func Is_Ok(Printable, Second_Printable string) {

	fmt.Println("> " + Printable + " Table was successfully created")

	if len(Second_Printable) > 0 {
		fmt.Println("-> " + Second_Printable + " was successfully created\n")
	} else {
		fmt.Println("")
	}

}
