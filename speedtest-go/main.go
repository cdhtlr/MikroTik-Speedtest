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
	buf		[]byte
	r		io.Reader
	iterNum		int
	startTime	time.Time
	// speeds
	avgSpd		float64
}

var (
	bufKB						= 50		// http buffer size in KB
	file			*os.File
	client			*http.Client
	result			string
	float_result, elapsed	float64
	err			error
	resp			*http.Response
	d			*downloader
	items			[]opts.LineData
	n			int
	fileScanner		*bufio.Scanner
	line			*charts.Line

)

func main() {
	createDB()
	
	http.HandleFunc("/", chart)
	
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {		
		start()
		fmt.Fprint(w, result)
	})
	
	http.HandleFunc("/condition", func(w http.ResponseWriter, r *http.Request) {		
		start()
		float_result, err = strconv.ParseFloat(result, 64)

		if err == nil {
			thresholdMbps, err := strconv.ParseFloat(os.Getenv("THRESHOLD_MBPS"), 64)
		
			if err == nil {
				// If the test results are above the threshold (Mbps) then the internet is in good condition
				if (float_result >= thresholdMbps){
					fmt.Fprint(w, "Good")
				} else {
					fmt.Fprint(w, "Bad")
				}
			} else {
				fmt.Println("Conversion Error: Cannot convert threshold Mbps to float data type")
				fmt.Fprint(w, "Unknown")
			}
		} else {
			fmt.Println("Conversion Error: Cannot convert download result to float data type")
			fmt.Fprint(w, "Unknown")
		}
	})
	
	fmt.Println("Speedtest web application runs on port 80")
	http.ListenAndServe(":80", nil)
}

func createDB(){
	// Create speedtest history file
	file, err = os.Create("data.txt")
	defer file.Close()

	if err != nil {
		fmt.Println("File Creation Error: Cannot create download history text file")
	}	
}

func start(){
	defer debug.FreeOSMemory()
	client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	
	// Download file from URL
	resp, err = client.Get(os.Getenv("URL"))
	defer resp.Body.Close()
	
	if err == nil {
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Invalid response ", resp.Status)
		}

		d = newDownloader(resp.Body)
		d.downSpeed()	
	} else {
		fmt.Println("Download Error: Cannot download "+os.Getenv("URL"))		
	}
}

func appendData(data string){
	file,err = os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	
	if err == nil {
		_, err = file.WriteString(data+"\n")

		if err != nil {
			fmt.Println("File Writing Error: Cannot write to download history text file")			
		}
	} else {
		fmt.Println("File Open Error: Cannot open download history text file")
	}
}

func generateChartItems() []opts.LineData {
	items = make([]opts.LineData, 0)
	
	file, err = os.Open("data.txt")
	defer file.Close()
	
	if err == nil {
		fileScanner = bufio.NewScanner(file)
		
		fileScanner.Split(bufio.ScanLines)
		
		for fileScanner.Scan() {
			items = append(items, opts.LineData{Value: fileScanner.Text()})
		}		
	} else {
		fmt.Println("File Open Error: Cannot open download history text file")
	}

	return items
}

func chart(w http.ResponseWriter, _ *http.Request) {
	defer debug.FreeOSMemory()
	line = charts.NewLine()
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
		buf:		make([]byte, 1024*bufKB),
		r:		r,
		startTime:	time.Now(),
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
		}
		
		// Stop speedtest after downloading MAX_KB
		maxKB, err := strconv.Atoi(os.Getenv("MAX_KB"))

		if err == nil {
			if d.iterNum*bufKB >= maxKB {
				break
			}		
		} else {
			fmt.Println("Conversion Error: Cannot convert max KB to integer data type")
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
