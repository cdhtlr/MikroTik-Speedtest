package main

import (
	"flag"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"io"
	"net/http"
	"time"
	
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type downloader struct {
	buf       []byte
	r         io.Reader
	iterNum   int
	startTime time.Time
	// speeds
	avgSpd float64
}

var (
	bufKB        	= 50    		// http buffer size in KB
	maxKB        	= 1000  		// stop speedtest after downloading maxKB
	thresholdMbps	= 1.0  			// If the test results are above the threshold (Mbps) then the internet is in good condition
	dataFile	= "data.txt"	// speedtest history file
	url *string
	client *http.Client
	result string
	float_result, elapsed float64
	err error
	_, file *os.File
	resp *http.Response
	d *downloader
	items []opts.LineData
	n int
	fileScanner *bufio.Scanner
	line *charts.Line
)

func main() {
	flag.IntVar(&maxKB, "m", 1000, "maximum size in KB to download")
	flag.Float64Var(&thresholdMbps, "t", 1.0, "download threshold in Mbps, to check for download speed condition")
	url = flag.String("u", "https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000", "url to download")

	flag.Parse()
	
	file, err = os.Create(dataFile)

	if err != nil {
		fmt.Println(err)
	}
	
	defer file.Close()
	
	http.HandleFunc("/", chart)
	
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {		
		start()
		fmt.Fprint(w, result)
	})
	
	http.HandleFunc("/condition", func(w http.ResponseWriter, r *http.Request) {		
		start()
		float_result, err = strconv.ParseFloat(result, 64)
		if err != nil {
			fmt.Println(err)
		}

		if (float_result >= thresholdMbps){
			fmt.Fprint(w, "Good")
		} else {
			fmt.Fprint(w, "Bad")
		}
	})
	
	fmt.Println("Speedtest web application runs on port 80")
	http.ListenAndServe(":80", nil)
}

func start(){
	client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	
	resp, err = client.Get(*url)
	
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Invalid response ", resp.Status)
	}

	d = newDownloader(resp.Body)
	d.downSpeed()
}

func appendData(data string){
	file,err = os.OpenFile(dataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	
	if err != nil {
		fmt.Println("Could not open speedtest history data")
		return
	}

	defer file.Close()
	 
	_, err = file.WriteString(data+"\n")

	if err != nil {
		fmt.Println("Could not write speedtest history data")
	}
}

func generateChartItems() []opts.LineData {
	items = make([]opts.LineData, 0)
	
	file, err = os.Open(dataFile)
	
	if err != nil {
		fmt.Println(err)
	}
	
	fileScanner = bufio.NewScanner(file)
	
	fileScanner.Split(bufio.ScanLines)
  
	for fileScanner.Scan() {
		items = append(items, opts.LineData{Value: fileScanner.Text()})
	}

    	file.Close()  
	return items
}

func chart(w http.ResponseWriter, _ *http.Request) {
	line = charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWalden, PageTitle: "Speedtest"}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Speedtest Chart",
			Subtitle: "Go to /test for speedtest or /condition to check for threshold based download speed condition",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true})
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
		buf:       make([]byte, 1024*bufKB),
		r:         r,
		startTime: time.Now(),
	}
}

func (d *downloader) downSpeed() {
	for {
		n, err = io.ReadFull(d.r, d.buf)
		_ = n
		d.iterNum++
		result = d.speedstr(true)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Println(err)
		}
		if d.iterNum*bufKB >= maxKB {
			break
		}
	}
	result = d.speedstr(false)
	appendData(result)
}

func (d *downloader) speeds() {
	elapsed = time.Since(d.startTime).Seconds()
	d.avgSpd = float64(d.iterNum*bufKB) / elapsed // in KB/s
}

func (d *downloader) speedstr(notFinalRun bool) string {
	if notFinalRun {
		d.speeds()
	}
	
	// converts download speed to Mb/s.
	return fmt.Sprintf("%.2f", d.avgSpd/1024*8)
}
