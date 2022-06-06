package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

//#------------------------------------------------------------------------------------------------------------# ↓ Select field for register or login ↓

//Multiple func Who_Whant (--> "Register" or "Login")
//Select field for log/register in bdd
func Check_If_Exist(input, input2, check_field, In_This_Table, Who_Want string) bool {

	var (
		db, _ = sql.Open(Bdd.Langage, Bdd.Name)
		Rows  []string
		rows  *sql.Rows
		I     = Instance_Bdd{}
		u     = all_bd{}
		index = &u.User.Name //<-- Default as Name & redefine on if
	)

	if Who_Want == "Register" { //<-- Define var by Who_want
		rows = Select_Field_From_DB(db, check_field, In_This_Table)
		input = strings.ToLower(input)
	} else if Who_Want == "Login" {
		rows = Select_All_From_DB(db, "user")
	}

	for rows.Next() {

		if Who_Want == "Register" { //<-- Select field
			if check_field == "Email" {
				index = &u.User.Email
			}
			err := rows.Scan(index)
			if err != nil {
				log.Fatal(err)
			}

			Rows = append(Rows, *index)

		} else if Who_Want == "Login" {

			err := rows.Scan(&u.User.Id, &u.User.Name, &u.User.Pswd, &u.User.Desc, &u.User.Email, &u.User.Profile_Picture, &u.User.Rank_id)
			if err != nil {
				log.Fatal(err)
			}

			u.User.Name = strings.ToLower(u.User.Name) //<-- Check in lower case to be sure
			u.User.Email = strings.ToLower(u.User.Email)
			I.I = append(I.I, u) //<-- Instance of User{}
		}

	}

	if Who_Want == "Register" { //<-- Select return
		return Check_Login_Or_Register(nil, "", "", Who_Want, input, Rows)
	} else if Who_Want == "Login" {
		return Check_Login_Or_Register(&I, input, input2, Who_Want, "", nil)
	}

	return false
}

//#------------------------------------------------------------------------------------------------------------# ↓ Check func for register or login ↓

//Check in dbb if exist
func Check_Login_Or_Register(I *Instance_Bdd, identifier, pswd, Who_whant, input string, Rows []string) bool {

	if Who_whant == "Login" {

		for _, i := range I.I { //<-- Check in I (= instance of table)
			if identifier == i.User.Email || identifier == i.User.Name {
				//<<<< Set for personnal Cookie
				Connected.User = i.User.Name
				Connected.User_Hased = Encrypt_Cookie(Connected.User)
				Connected.Rank_Id = strconv.Itoa(i.User.Rank_id)
				Connected.Rank_Id_Hashed = Encrypt_Cookie(strconv.Itoa(i.User.Rank_id))
				//<<<< Set for personnal Cookie
				return CheckPasswordHash(pswd, i.User.Pswd)
			}
		}
		return false

	} else if Who_whant == "Register" {

		for i := range Rows {
			low := strings.ToLower(Rows[i]) //<-- Check in lower case to be sure
			if low == input {
				return false
			}
		}

		return true
	}
	return false
}
