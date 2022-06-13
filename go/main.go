package main

import (
	"fmt"
	Database "forum/database"

	Config "forum/config"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	Database.Terminal_Init_Table("Bdd_Name")
	Database.InitBDD()
	fmt.Print(Config.Bdd.Name)
	HttpServ()
}
