package mrcong

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"page.github.io/pkg/file"
	"path"
	"path/filepath"
	"strings"
)

type jsonData struct {
	MediafireLink []string
	Title         string
	Url           string
}

func DownloadMediafireLink() {
	dirName := "json/" + file.GetNameWithoutExt()
	var files []string
	extName := ".json"
	err := filepath.Walk(dirName, func(pathFile string, info os.FileInfo, err error) error {
		if extName != path.Ext(pathFile) {
			return nil
		}

		content, err := os.ReadFile(pathFile)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Now let's unmarshall the data into `payload`
		var payload jsonData
		err = json.Unmarshal(content, &payload)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if len(payload.MediafireLink) == 0 {
			return nil
		}

		newPageLinks := FilterInvalidMediafireLink(payload.MediafireLink)
		if len(newPageLinks) == 0 {
			return nil
		}

		for _, pageLink := range newPageLinks {
			// Request the HTML page.
			res, err := http.Get(pageLink)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(res.Body)
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}

			downloadLink, _ := doc.Find("#downloadButton").Attr("href")
			fmt.Println(pathFile)
			fmt.Println(payload)
			fmt.Println(pageLink)
			fmt.Println(downloadLink)

			fileURL, err := url.Parse(downloadLink)
			if err != nil {
				fmt.Println(err)
				continue
			}
			segments := strings.Split(fileURL.Path, "/")
			fileName := segments[len(segments)-1]
			outputFile := "images/" + file.GetNameWithoutExt() + "/" + fileName

			if file.Exists(outputFile) {
				continue
			}

			err = got.New().Download(downloadLink, outputFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		files = append(files, pathFile)
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}

func FilterInvalidMediafireLink(oldLinks []string) (newLinks []string) {
	validLinks := []string{"https://www.mediafire.com"}
	//invalidLinks := []string{"http://shink.me", "http://ouo.io", "http://adf.ly"}

	for _, value := range oldLinks {
		for _, link := range validLinks {
			if strings.Contains(value, link) {
				newLinks = append(newLinks, value)
				break
			}
		}

		//for _, link := range invalidLinks {
		//	if strings.Contains(value, link) {
		//		continue level1
		//	}
		//}

	}

	return
}
