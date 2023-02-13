package download

import (
	"bufio"
	"crypto/tls"
	. "io"
	"math"
	"net/http"
	"os"
	. "strconv"
	"strings"
	. "time"
)

type buf_downloader struct {
	iterNum int
	buf     []byte
	r       Reader
}

var (
	downloadCompleted = 0
	max_dl_size_int   = int(math.Floor(max_dl_size)) * 1024 * 1024
	max_dl_size, _    = ParseFloat(os.Getenv("MAX_DLSIZE"), 64)
	threshold, _      = ParseFloat(os.Getenv("MIN_THRESHOLD"), 64)
	url               = os.Getenv("URL")
)

func isAcceptRangeSupported() (bool, int) {
	req, _ := http.NewRequest("HEAD", url, nil)
	client := &http.Client{
		Timeout: 5 * Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return false, 0
	}

	acceptRanges := strings.ToLower(resp.Header.Get("Accept-Ranges"))
	if acceptRanges == "" || acceptRanges == "none" {
		return false, int(resp.ContentLength)
	}

	return true, int(resp.ContentLength)
}

func downloadPart(start int, end int, done chan bool) {
	download("NO", Itoa(int(start))+"-"+Itoa(int(end)))
	done <- true
}

func download(opts ...string) {
	req, _ := http.NewRequest("GET", url, nil)
	if len(opts) > 1 {
		req.Header.Add("Range", "bytes="+opts[1])
	}
	client := &http.Client{
		Timeout: 30 * Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		downloadCompleted -= 1
	} else {
		downloadCompleted += 1
	}
	defer resp.Body.Close()

	if opts[0] == "YES" {
		d := &buf_downloader{
			buf: make([]byte, 1024),
			r:   resp.Body,
		}
		d.bufDown()
	} else {
		Copy(Discard, resp.Body)
	}
}

func (d *buf_downloader) bufDown() {
	for {
		_, err := ReadFull(d.r, d.buf)
		d.iterNum++
		if err == EOF || err == ErrUnexpectedEOF || float64((d.iterNum/1024)) >= float64(max_dl_size_int/1024/1024) {
			break
		}
	}
}

func saveResult(speedtest_result string) error {
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(speedtest_result + "\n")

	if err != nil {
		return err
	}

	return writer.Flush()
}

func Run(explain bool) (int, string) {
	downloadCompleted = 0

	concurentConn, _ := Atoi(os.Getenv("CONCURENT_CONNECTION"))

	acceptRangeSupported, fileSize := isAcceptRangeSupported()

	begin := Now()

	if fileSize > 0 {
		if 0 < max_dl_size_int && max_dl_size_int < fileSize {
			fileSize = max_dl_size_int
		}

		if acceptRangeSupported {
			partSize := fileSize / concurentConn
			done := make(chan bool, concurentConn)

			for i := 0; i < concurentConn; i++ {
				start := i * partSize
				end := (i+1)*partSize - 1
				if i == concurentConn-1 {
					end = fileSize - 1
				}
				go downloadPart(start, end, done)
			}
			for i := 0; i < concurentConn; i++ {
				<-done
			}
		} else {
			concurentConn = 1
			download(strings.ToUpper(os.Getenv("ALLOW_MEMORY_BUFFER")))
		}
	}

	elapsed_time := Since(begin).Seconds()
	speedtest_result := float64(fileSize*8/1024/1024) / elapsed_time
	speedtest_result_string := FormatFloat(speedtest_result, 'f', 2, 64)
	code := 200
	condition := "Good"

	if speedtest_result < threshold {
		code = 201
		condition = "Bad"
	}

	if speedtest_result_string == "NaN" || downloadCompleted != concurentConn {
		speedtest_result_string = "0.0"
	}

	saveResult(speedtest_result_string)

	if explain {
		return code, condition
	}

	return code, speedtest_result_string
}
