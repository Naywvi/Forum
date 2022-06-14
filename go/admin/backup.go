package admin

import (
	"encoding/json"
	Config "forum/config"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func Set_Backup(I_I Config.Instance_of_instance) {
	var (
		count        = 0
		number_files = strconv.Itoa(count)
		files, err   = ioutil.ReadDir("../backup/")
		t            = time.Now()
		date         = t.Format("2006-01-02")
		backup, _    = json.Marshal(I_I)
	)
	if err != nil {
		log.Fatal(err)
	}
	for range files {
		count++
	}
	if count > 0 {
		number_files = "(" + strconv.Itoa(count) + ")"
	} else {
		number_files = ""
	}

	f, err := os.Create("../backup/Backup_of_" + Config.Bdd.Name + "_" + date + number_files + ".json")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(backup)
}
