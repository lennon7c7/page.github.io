package mrcong

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"io"
	"net/http"
	"net/url"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/log"
	"page.github.io/pkg/proxy"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var IsExit bool
var Domain = "https://mrcong.com"
var BaseDownloadJsonPath = "../../json/" + file.GetNameWithoutExt() + "/"
var BaseDownloadZipPath = "../../zip/" + file.GetNameWithoutExt() + "/"

type jsonData struct {
	MediafireLink []string
	Title         string
	Updated       string
	Url           string
	ImgLinkList   []string
}

func DownloadMediafireLink(jsonFiles []string) {
	if len(jsonFiles) == 0 {
		extName := ".json"
		err := filepath.Walk(BaseDownloadJsonPath, func(pathFile string, info os.FileInfo, err error) error {
			if extName != path.Ext(pathFile) {
				return nil
			}

			jsonFiles = append(jsonFiles, pathFile)

			return nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		// 降序
		sort.Sort(sort.Reverse(sort.StringSlice(jsonFiles)))
	}

	for _, jsonFile := range jsonFiles {
		jsonFile, err := filepath.Abs(jsonFile)
		if err != nil {
			log.Error(err)
			continue
		}

		content, err := os.ReadFile(jsonFile)
		if err != nil {
			log.Error(err)
			continue
		}

		// Now let's unmarshall the data into `payload`
		var payload jsonData
		err = json.Unmarshal(content, &payload)
		if err != nil {
			log.Error(err)
			continue
		}

		if len(payload.MediafireLink) == 0 {
			continue
		}

		newPageLinks := FilterInvalidMediafireLink(payload.MediafireLink)
		if len(newPageLinks) == 0 {
			continue
		}

		for _, pageLink := range newPageLinks {
			// Request the HTML page.
			res, err := http.Get(pageLink)
			if err != nil {
				log.Error(err)
				continue
			}
			defer func(Body io.ReadCloser) {
				closeErr := Body.Close()
				if closeErr != nil {
					err = closeErr
					log.Error(err)
				}
			}(res.Body)
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
				continue
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Error(err)
				continue
			}

			downloadLink, _ := doc.Find("#downloadButton").Attr("href")

			fileURL, err := url.Parse(downloadLink)
			if err != nil {
				log.Error(err)
				continue
			}
			segments := strings.Split(fileURL.Path, "/")
			fileName := segments[len(segments)-1]
			outputFile := BaseDownloadZipPath + fileName
			outputFile, err = filepath.Abs(outputFile)
			if err != nil {
				log.Error(err)
				continue
			}

			if file.Exists(outputFile) {
				continue
			}

			fmt.Printf("jsonFile: %v\n", jsonFile)
			err = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("  ", downloadLink)
			err = got.New().Download(downloadLink, outputFile)
			if err != nil {
				log.Error(err)
				continue
			}

			fmt.Println("  ", outputFile)
		}
	}
}

func DownloadToJson(webUrl string) {
	listPage(webUrl)
}

func listPage(url string) {
	IsExit = false
	fmt.Println(url)

	doc, err := proxy.GetHtmlDom(url)
	if err != nil {
		log.Error(err)
		return
	}

	doc.Find(".post-listing .post-box-title a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		if detailUrl != "" {
			err = detailPage(detailUrl)
			if err != nil {
				log.Error(err)
				return
			}
		}
	})

	if IsExit {
		return
	}

	nextUrl, _ := doc.Find("head link[rel=next]").Attr("href")
	if nextUrl != "" {
		listPage(nextUrl)
	}
}

func detailPage(url string) (err error) {
	if url == "" {
		err = errors.New(`url == ""`)
		return
	}
	fmt.Println("  ", url)

	doc, err := proxy.GetHtmlDom(url)
	if err != nil {
		return
	}

	title := doc.Find("#crumbs .current").Text()
	if title == "" {
		err = errors.New(`title == ""`)
		return
	}

	downloadPath := doc.Find(".updated").Text()
	downloadPath = strings.Replace(downloadPath, "-", "/", 2)
	downloadPath = BaseDownloadJsonPath + downloadPath + "/"

	if !file.Exists(downloadPath) {
		err = os.MkdirAll(downloadPath, os.ModePerm)
		if err != nil {
			return
		}
	}

	jsonFile := downloadPath + filepath.Base(url) + ".json"
	jsonFile, err = filepath.Abs(jsonFile)
	if err != nil {
		return
	}

	if file.Exists(jsonFile) {
		fmt.Println("---------- no shit ---------- ")
		IsExit = true
		return
	}

	var mediafireLinkList []string
	doc.Find("a.shortc-button.medium.green").Each(func(i int, s *goquery.Selection) {
		mediafireLink, _ := s.Attr("href")
		mediafireLinkList = append(mediafireLinkList, mediafireLink)
	})

	var imgLinkList []string
	imgLinkList = getDetailPageImgList(url)
	if len(imgLinkList) == 0 {
		err = errors.New(`len(imgLinkList) == 0`)
		return
	}

	dataMap, err := json.MarshalIndent(jsonData{
		MediafireLink: mediafireLinkList,
		Title:         title,
		Updated:       doc.Find(".updated").Text(),
		Url:           url,
		ImgLinkList:   imgLinkList,
	}, "", "  ")
	if err != nil {
		return
	}

	err = os.WriteFile(jsonFile, dataMap, os.ModePerm)
	if err != nil {
		return
	}

	fmt.Println("    ", jsonFile)

	return
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

func getDetailPageImgList(detailPage string) (imgList []string) {
	//fmt.Println("    ", detailPage)

	doc, err := proxy.GetHtmlDom(detailPage)
	if err != nil {
		log.Error(err)
		return
	}

	doc.Find("#fukie2 img.aligncenter").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("src")
		imgList = append(imgList, val)
		//fmt.Println("      ", val)
	})

	nextDetailPage, _ := doc.Find("head link[rel=next]").Attr("href")
	if nextDetailPage != "" {
		temp := getDetailPageImgList(nextDetailPage)
		imgList = append(imgList, temp...)
	}

	return
}
