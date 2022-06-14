package config

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Connected_Status struct {
	User           string
	User_Hased     string
	Rank_Id        string
	Rank_Id_Hashed string
}

var Connected Connected_Status
var TemplatesDir = os.Getenv("TEMPLATES_DIR")

//#------------------------------------------------------------------------------------------------------------# ↓ Return to [Select_Page] ↓

//Return to page Selected need Path (string)
func Return_To_Page(w http.ResponseWriter, r *http.Request, Path string) {
	template.Must(template.ParseFiles(filepath.Join(TemplatesDir, Path))).Execute(w, " ")
}

//#------------------------------------------------------------------------------------------------------------# ↓ Return error html ↓

//Send Http error method
func Send_Error(w http.ResponseWriter, r *http.Request) {
	Return_To_Page(w, r, "../static/templates/managed_pages/404.html")
	// http.Error(w, "Method is not supported.", http.StatusBadRequest) //<-- Print [error] Method is not supported
	// fmt.Fprint(w, http.StatusBadRequest)

}

func HandleError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
