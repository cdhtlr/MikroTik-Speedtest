package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/showwin/speedtest-go/speedtest"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	showList   = kingpin.Flag("list", "Show available speedtest.net servers.").Short('l').Bool()
	serverIds  = kingpin.Flag("server", "Select server id to speedtest.").Short('s').Ints()
	uploadTest = kingpin.Flag("upload", "Upload test.").Short('u').Bool()
	dwloadTest = kingpin.Flag("download", "Download test.").Short('d').Bool()
	savingMode = kingpin.Flag("saving-mode", "Using less memory (≒10MB), though low accuracy (especially > 30Mbps).").Bool()
	jsonOutput = kingpin.Flag("json", "Output results in json format").Bool()
)

type fullOutput struct {
	Timestamp outputTime        `json:"timestamp"`
	UserInfo  *speedtest.User   `json:"user_info"`
	Servers   speedtest.Servers `json:"servers"`
}
type outputTime time.Time

func main() {
	str := make(chan string, 1)
    go func() {
		kingpin.Version("1.1.5")
		kingpin.Parse()

		user, err := speedtest.FetchUserInfo()
		if err != nil {
			fmt.Println("Warning: Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
		}
		if !*jsonOutput {
			showUser(user)
		}

		servers, err := speedtest.FetchServers(user)
		checkError(err)
		if *showList {
			showServerList(servers)
			return
		}

		targets, err := servers.FindServer(*serverIds)
		checkError(err)

		startTest(targets, *uploadTest, *dwloadTest, *savingMode, *jsonOutput)

		if *jsonOutput {
			jsonBytes, err := json.Marshal(
				fullOutput{
					Timestamp: outputTime(time.Now()),
					UserInfo:  user,
					Servers:   targets,
				},
			)
			checkError(err)

			fmt.Println(string(jsonBytes))
		}
        str <- "ok"
    }()
    select {
		case <-time.After(61 * time.Second):
			os.Exit(1)
    }
}

func startTest(servers speedtest.Servers, uploadTest bool, dwloadTest bool, savingMode bool, jsonOutput bool) {
	for _, s := range servers {
		if !jsonOutput {
			showServer(s)
		}

		err := s.PingTest()
		checkError(err)

		if jsonOutput {
			if uploadTest {
				err = s.UploadTest(savingMode)
				checkError(err)
			}
			if dwloadTest {
				err := s.DownloadTest(savingMode)
				checkError(err)
			}

			continue
		}

		showLatencyResult(s)

		if uploadTest {
			err = testUpload(s, savingMode)
			checkError(err)
		}

		if dwloadTest {
			err = testDownload(s, savingMode)
			checkError(err)
		}
		
		showServerResult(s)
	}

	if !jsonOutput && len(servers) > 1 {
		showAverageServerResult(servers)
	}
}

func testUpload(server *speedtest.Server, savingMode bool) error {
	quit := make(chan bool)
	go waiting(quit)
	err := server.UploadTest(savingMode)
	quit <- true
	if err != nil {
		return err
	}
	return err
}

func testDownload(server *speedtest.Server, savingMode bool) error {
	quit := make(chan bool)
	go waiting(quit)
	err := server.DownloadTest(savingMode)
	quit <- true
	if err != nil {
		return err
	}
	return err
}

func waiting(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

func showUser(user *speedtest.User) {
	if user.IP != "" {
		fmt.Printf("Testing From IP: %s\n", user.String())
	}
}

func showServerList(servers speedtest.Servers) {
	for _, s := range servers {
		fmt.Printf("[%4s] %8.2fkm ", s.ID, s.Distance)
		fmt.Printf(s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
	}
}

func showServer(s *speedtest.Server) {
	fmt.Printf("Target Server: [%4s] %8.2fkm ", s.ID, s.Distance)
	fmt.Printf(s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
}

func showLatencyResult(server *speedtest.Server) {
	fmt.Println("Latency:", server.Latency)
}

// ShowResult : show testing result
func showServerResult(server *speedtest.Server) {
	fmt.Printf("Upload: %5.2f Mbit/s\n", server.ULSpeed)
	fmt.Printf("Download: %5.2f Mbit/s\n", server.DLSpeed)
	valid := server.CheckResultValid()
	if !valid {
		fmt.Println("Warning: Result seems to be wrong. Please speedtest again.")
	}
}

func showAverageServerResult(servers speedtest.Servers) {
	avgUL := 0.0
	avgDL := 0.0
	for _, s := range servers {
		if (s.ULSpeed>0){
			avgUL = avgUL + s.ULSpeed
		}
		if (s.DLSpeed>0){
			avgDL = avgDL + s.DLSpeed
		}
	}
	
	if (avgUL>0){
		fmt.Printf("Upload Avg: %5.2f Mbit/s\n", avgUL/float64(len(servers)))
	}
	if (avgDL>0){
		fmt.Printf("Download Avg: %5.2f Mbit/s\n", avgDL/float64(len(servers)))
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (t outputTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05.000"))
	return []byte(stamp), nil
}
