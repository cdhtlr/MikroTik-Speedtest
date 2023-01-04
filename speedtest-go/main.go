package main

import (
	"os"
	"fmt"
	"net/http"
	"strconv"
	"runtime/debug"
	"crypto/tls"
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
	result		float64
)

func main() {
	createDB()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {		
		start()
		
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(resultToCode(result))

		fmt.Fprintf(w, "%.2f", result)
	})
	
	http.HandleFunc("/condition", func(w http.ResponseWriter, r *http.Request) {
		start()
		
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.WriteHeader(resultToCode(result))

		if (resultToCode(result) == 200){
			fmt.Fprint(w, "Good")
		} else {
			fmt.Fprint(w, "Bad")			
		}
	})
	
	http.HandleFunc("/chart", chart)

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
			MaxIdleConnsPerHost: 1,
			IdleConnTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
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

func resultToCode(result float64) int {
	threshold, _ := strconv.ParseFloat(os.Getenv("MIN_THRESHOLD"), 64)
	
	// If the test result is above the threshold (Mbps) or equal then the internet is in good condition	
	if (result >= threshold){
		return 200
	}
	
	return 201	
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

	items := generateChartItems()
	
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWalden, PageTitle: "Speedtest Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)

	line.SetXAxis(items).
		AddSeries("Download Speed", items).
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
		
		// Stop speedtest after downloading MAX_DLSIZE (in MB)
		maxDLSIZE, _ := strconv.ParseFloat(os.Getenv("MAX_DLSIZE"), 64)
		
		if float64(d.iterNum*50/1024) >= maxDLSIZE {
			break
		}
	}
	
	result = d.speedres(false)
	appendData(result)
}

func (d *downloader) speeds() {
	elapsed := time.Since(d.startTime).Seconds()
	d.avgSpd = float64(d.iterNum*50/1024*8) / elapsed // in Mb/s
}

func (d *downloader) speedres(notFinalRun bool) float64 {
	if notFinalRun {
		d.speeds()
	}
	
	return d.avgSpd
}
