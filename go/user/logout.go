package user

import (
	"fmt"
	Database "forum/database"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	var (
		_, _, User = Check_Cookie(w, r)
		logout     = time.Now()
	)

	Database.Update_Field("profil", "Last_time_connected", "User", User, logout.String())
	del(w, r)

	fmt.Fprint(w, `<script language="javascript" type="text/javascript"> window.location="/forum"; </script>`)
}
