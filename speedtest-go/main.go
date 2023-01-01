package main

import (
	"os"
	"fmt"
	"net/http"
	"strconv"
	"runtime/debug"
	"bufio"
	"io"
	"time"
	
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type downloader struct {
	startTime	time.Time
	avgSpd		float64
	iterNum		int
	buf		[]byte
	r		io.Reader
}

var (
	result			float64
)

func main() {
	createDB()
	
	http.HandleFunc("/", chart)
	
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {		
		start()
		fmt.Fprintf(w, "%.2f", result)
	})
	
	http.HandleFunc("/condition", func(w http.ResponseWriter, r *http.Request) {
		start()

		threshold, _ := strconv.ParseFloat(os.Getenv("THRESHOLD"), 64)
		
		// If the test results are above the threshold (Mbps) then the internet is in good condition
		if (result >= threshold){
			fmt.Fprint(w, "Good")
		} else {
			fmt.Fprint(w, "Bad")
		}
	})
	
	fmt.Println("Speedtest web application runs on port 80")
	http.ListenAndServe(":80", nil)
}

func createDB(){
	// Create speedtest history file
	file, _ := os.Create("data.txt")
	defer file.Close()
}

func start(){
	defer debug.FreeOSMemory()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	
	// Download file from URL
	resp, _ := client.Get(os.Getenv("URL"))
	defer resp.Body.Close()
	
	defer client.CloseIdleConnections()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Invalid response ", resp.Status)
	}

	d := newDownloader(resp.Body)
	d.downSpeed()
}

func appendData(data float64){
	file, _ := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	
	file.WriteString(fmt.Sprintf("%.2f", data)+"\n")
}

func generateChartItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	
	file, _ := os.Open("data.txt")
	defer file.Close()
	
	fileScanner := bufio.NewScanner(file)
	
	fileScanner.Split(bufio.ScanLines)
	
	for fileScanner.Scan() {
		items = append(items, opts.LineData{Value: fileScanner.Text()})
	}

	return items
}

func chart(w http.ResponseWriter, _ *http.Request) {
	defer debug.FreeOSMemory()

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWalden, PageTitle: "Speedtest"}),
		charts.WithTitleOpts(opts.Title{
			Title:		"Speedtest Chart",
			Subtitle:	"Go to /test for speedtest or /condition to check for threshold based download speed condition",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)

	line.SetXAxis(generateChartItems()).
		AddSeries("Download Speed", generateChartItems()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}),charts.WithMarkPointNameTypeItemOpts(
			opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
			opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
		))
	line.Render(w)
}

func newDownloader(r io.Reader) *downloader {
	return &downloader{
		// http buffer size in KB is 50
		buf:		make([]byte, 1024*50),
		r:		r,
		startTime:	time.Now(),
	}
}

func (d *downloader) downSpeed() {
	for {
		n, err := io.ReadFull(d.r, d.buf)
		_ = n
		d.iterNum++
		result = d.speedres(true)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
		}
		
		// Stop speedtest after downloading MAX_KB
		maxKB, _ := strconv.Atoi(os.Getenv("MAX_KB"))
		
		if d.iterNum*50 >= maxKB {
			break
		}
	}
	
	result = d.speedres(false)
	appendData(result)
}

func (d *downloader) speeds() {
	elapsed := time.Since(d.startTime).Seconds()
	d.avgSpd = float64(d.iterNum*50) / elapsed // in KB/s
}

func (d *downloader) speedres(notFinalRun bool) float64 {
	if notFinalRun {
		d.speeds()
	}
	
	// converts download speed to Mb/s.
	return d.avgSpd/1024*8
}
