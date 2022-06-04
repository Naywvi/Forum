package main

import (
	_ "github.com/mattn/go-sqlite3"
)

type Bdds struct {
	Name    string
	Langage string
}

var Bdd Bdds

func main() {
	Terminal_Init_Table("Bdd_Name")
	InitBDD()
	httpServ()
}
