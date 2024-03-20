package process_dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Callback interface {
	GetDiv(canvasId string) string
	GetJS(canvasId string) string
}

type DashBoard struct {
	Task      func() map[string]string
	Iters     int
	Callbacks []Callback
}

func (d DashBoard) Run() {

	data := make([]map[string]string, 0)

	mainPageHandler := func(w http.ResponseWriter, r *http.Request) {
		content := GetIndexContent()
		replaceCode := ""
		replaceCanvas := ""
		for id, code := range d.Callbacks {
			replaceCode += "\n"
			replaceCode += code.GetJS(fmt.Sprintf("%d", id))
			replaceCanvas += code.GetDiv(fmt.Sprintf("%d", id))
		}
		content = strings.Replace(content, "[CANVAS]", replaceCanvas, -1)
		content = strings.Replace(content, "[CODE]", replaceCode, -1)
		fmt.Fprintf(w, content)
	}

	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		j, _ := json.Marshal(data)
		fmt.Fprintf(w, strings.Replace(string(j), "[DATA]", string(j), 1))
	}

	go func() {
		for i := 0; i < d.Iters; i++ {
			data = append(data, d.Task())
		}
	}()

	http.HandleFunc("/data", apiHandler)
	http.HandleFunc("/", mainPageHandler)
	http.ListenAndServe(":8080", nil)

}

type LineCallback struct {
	Title    string
	Variable string
}

func (c LineCallback) GetDiv(id string) string {
	return fmt.Sprintf(`<div class="card">
		<div class="title">%s</div>
		<canvas id='chart-%s' class='chart'></canvas>
	</div>`, c.Title, id)
}

func (c LineCallback) GetJS(id string) string {
	return fmt.Sprintf(`
		new Chart("chart-%s", {
			type: "line",
			data: {
				labels: data.map((value, id) => id),
				datasets: [{
					borderColor: "rgba(0,232,98, 1)",
					data: data.map((value, id) => parseFloat(value["%s"])),
					fill: false
				}]
			}, options: {
				legend: {display: false},
				plugins: {
					title: {
						display: true,
						text: '%s'
					}
				}, animation: {
					duration: 0
				}
			}
		});
	`, id, c.Variable, c.Title)
}

type LastValueCallback struct {
	Title    string
	Variable string
}

func (c LastValueCallback) GetDiv(id string) string {
	return fmt.Sprintf(`<div class="card">
		<div class="title">%s</div>
		<div id="text-%s"></div>
	</div>`, c.Title, id)
}

func (c LastValueCallback) GetJS(id string) string {
	return fmt.Sprintf(`
		document.getElementById("text-%s").innerHTML = data[data.length-1]["%s"]
	`, id, c.Variable)
}
