package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ghoulseek/downloader/slsk"
	"ghoulseek/globals"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// TODO: implement error handling/returns

// this assumes all files are hosted by the same user
func StartDownload(files []slsk.File, userName string) {
	fmt.Println(files)
	fmt.Println(userName)
	urlToUse := globals.SlskdEndpoint + "/transfers/downloads/" + userName

	type DownloadParam struct {
		FileName string `json:"filename"`
		Size     int64  `json:"size"`
	}

	paramList := []DownloadParam{}

	for _, f := range files {
		paramList = append(paramList, DownloadParam{
			FileName: f.FileName,
			Size:     int64(f.Size),
		})
	}

	encode, _ := json.Marshal(paramList)

	resp, err := http.Post(urlToUse, "application/json", bytes.NewReader(encode))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	fmt.Println("Start download: " + strconv.Itoa(resp.StatusCode))

	c, _ := io.ReadAll(resp.Body)

	fmt.Println(string(c))
}

type DownloadFile struct {
	Id               string  `json:"id"`
	UserName         string  `json:"username"`
	Direction        string  `json:"direction"`
	FileName         string  `json:"filename"`
	Size             int64   `json:"size"`
	State            string  `json:"state"`
	BytesTransferred int64   `json:"bytesTransferred"`
	BytesRemaining   int64   `json:"bytesRemaining"`
	PercentComplete  float64 `json:"percentComplete"`
	AverageSpeed     float64 `json:"averageSpeed"`
}

type DownloadDirectory struct {
	Path      string         `json:"directory"`
	FileCount int            `json:"fileCount"`
	Files     []DownloadFile `json:"files"`
}

type DownloadEntry struct {
	UserName    string              `json:"username"`
	Directories []DownloadDirectory `json:"directories"`
}

// todo: implement error handling and returns
func GetDownloadList() []DownloadEntry {
	urlToUse := globals.SlskdEndpoint + "/transfers/downloads/"

	resp, err := http.Get(urlToUse)
	if err != nil {
		fmt.Println(err.Error())
		return []DownloadEntry{}
	}
	defer resp.Body.Close()

	var result []DownloadEntry
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return []DownloadEntry{}
	}

	return result
}

func WaitForDownloadFinish(files []slsk.File, userName string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		var hasIncompleteFile bool = false

		select {
		case <-ticker.C:
			downloads := GetDownloadList()

			// holy recursion
			for _, d := range downloads {
				if hasIncompleteFile {
					break
				}
				if d.UserName == userName {
					for _, dir := range d.Directories {
						for _, f := range dir.Files {
							for _, wf := range files {
								if wf.FileName == f.FileName {
									if f.BytesRemaining != 0 {
										fmt.Println("File " + color.GreenString(f.FileName) + " is incomplete")
										hasIncompleteFile = true
									}
								}
							}
						}
					}
					break
				}
			}
			if !hasIncompleteFile {
				fmt.Println("download finished")
				return
			}
		}
	}
}
