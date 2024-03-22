package process_dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
		start1 := time.Now()
		for i := 0; i < d.Iters; i++ {
			start := time.Now()
			taskResult := d.Task()
			elapsed := time.Since(start)
			taskResult["currentIteration"] = fmt.Sprintf("%d", i)
			taskResult["elapsed_time"] = fmt.Sprintf("%s", elapsed)
			data = append(data, taskResult)
		}
		fmt.Println("Work done !")
		fmt.Println(time.Since(start1))
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
					borderColor: "#17A589",
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

type MultiLineCallback struct {
	Title          string
	Variables      []string
	VariablesNames []string
}

func (c MultiLineCallback) GetDiv(id string) string {
	return fmt.Sprintf(`<div class="card">
		<div class="title">%s</div>
		<canvas id='chart-%s' class='chart'></canvas>
	</div>`, c.Title, id)
}

func (c MultiLineCallback) GetJS(id string) string {
	colors := []string{"#17A589", "#2471A3", "#7D3C98", "#CB4335", "#F1C40F"}
	text := fmt.Sprintf(`
		new Chart("chart-%s", {
			type: "line",
			data: {
				labels: data.map((value, id) => id),
				datasets: [`, id)
	for id, s := range c.Variables {
		text += fmt.Sprintf(`{
			data: data.map((value, id) => parseFloat(value["%s"])),
			fill: false,
			label: "%s",
			borderColor: "%s",
		},`, s, c.VariablesNames[id], colors[id%len(colors)])
	}
	text += fmt.Sprintf(`]
		}, options: {
			plugins: {
				title: {
					display: true,
					text: '%s'
				}
			}, animation: {
				duration: 0
			}
		}
	});`, c.Title)
	return text
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
