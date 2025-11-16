package controller

import (
	"html/template"
	"net/http"

	"TP_Andreev/internal/transport/http/router"
)

var tpl = template.Must(template.ParseFiles("web/templates/index.html"))

type MainController struct{}

// GetMainPage handles GET / requests
func (pc *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	tpl.Execute(w, nil)
}
