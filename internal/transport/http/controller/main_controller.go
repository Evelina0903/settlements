package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"settlements/internal/service"
	"settlements/internal/transport/http/router"
)

type MainController struct {
	service service.Service
}

type tmplData struct {
	Table  template.JS
	Chart1 template.JS
	Chart2 template.JS
}

var tmpl = template.Must(
	template.Must(
		template.New("jsData").Parse(src),
	).ParseFiles("web/templates/index.html"),
)

func (c *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	settelmentType := c.service.GetAllSettelmetTypeData()
	settelmentTypeJ, _ := json.Marshal(settelmentType)

	longitudePopulation := c.service.GetLongitudePopulationData()
	longitudePopulationJ, _ := json.Marshal(longitudePopulation)

	districtPopulation := c.service.GetDistrictPopulationData()
	districtPopulationJ, _ := json.Marshal(districtPopulation)

	data := tmplData{
		Table:  template.JS(settelmentTypeJ),
		Chart1: template.JS(longitudePopulationJ),
		Chart2: template.JS(districtPopulationJ),
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

const src = `
	<script>
        const tableData = {{.Table}};
        const chartData1 = {{.Chart1}};
        const chartData2 = {{.Chart2}};
    </script>`
