package main

import (
	. "bufio"
	"net/http"
	. "os"
	"runtime/debug"

	"MikroTik-Speedtest/function/download"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	createDB()

	http.HandleFunc("/", speedtestHandler)
	http.HandleFunc("/condition", speedtestHandler)
	http.HandleFunc("/chart", chart)

	Stdout.Write([]byte("Speedtest web application runs on port 80\n"))
	http.ListenAndServe(":80", nil)
}

func createDB() {
	file, _ := Create("data.txt")
	defer file.Close()
}

func speedtestHandler(w http.ResponseWriter, r *http.Request) {
	defer debug.FreeOSMemory()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	code := 404
	result := "Not found"

	if r.URL.Path == "/" {
		code, result = download.Run(false)
	} else if r.URL.Path == "/condition" {
		code, result = download.Run(true)
	}
	w.WriteHeader(code)
	w.Write([]byte(result))
}

func generateChartItems() []opts.LineData {
	file, _ := Open("data.txt")
	defer file.Close()

	fileScanner := NewScanner(file)
	fileScanner.Split(ScanLines)

	items := make([]opts.LineData, 0)
	for fileScanner.Scan() {
		items = append(items, opts.LineData{Value: fileScanner.Text()})
	}
	return items
}

func chart(w http.ResponseWriter, _ *http.Request) {
	defer debug.FreeOSMemory()

	items := generateChartItems()

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWalden, PageTitle: "Speedtest Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)

	line.SetXAxis(items).
		AddSeries("Download Speed", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}), charts.WithMarkPointNameTypeItemOpts(
			opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
			opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
		))
	line.Render(w)
}
