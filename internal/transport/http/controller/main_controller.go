package controller

import (
	"html/template"
	"net/http"

	"settlements/internal/transport/http/router"
)

var mainTpl = template.Must(template.ParseFiles("web/templates/index.html"))

type MainController struct{}

// GetMainPage handles GET / requests
func (pc *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	mainTpl.Execute(w, nil)
}
