package mrcong

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"page.github.io/pkg/file"
	"path"
	"path/filepath"
	"strings"
)

var domain = "https://mrcong.com"
var BaseDownloadJsonPath = "../../json/" + file.GetNameWithoutExt() + "/"

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

		content, err := ioutil.ReadFile(pathFile)
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
			outputFile := "zip/" + file.GetNameWithoutExt() + "/" + fileName

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

func DownloadToJson() {
	listPage(domain)
}

func listPage(url string) {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".post-listing .post-box-title a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		if detailUrl != "" {
			detailPage(detailUrl, 0)
		}
	})

	nextUrl, _ := doc.Find("head link[rel=next]").Attr("href")
	if nextUrl != "" {
		listPage(nextUrl)
	}
}

func detailPage(url string, page int) {
	fmt.Println("  " + url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find("#crumbs .current").Text()

	downloadPath := doc.Find(".updated").Text()
	downloadPath = strings.Replace(downloadPath, "-", "/", 2)
	downloadPath = BaseDownloadJsonPath + downloadPath + "/"

	if !file.Exists(downloadPath) {
		err = os.MkdirAll(downloadPath, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	jsonFile := downloadPath + title + ".json"
	fmt.Println("  " + jsonFile)
	if file.Exists(jsonFile) {
		fmt.Println("---------- no shit ---------- ")
		os.Exit(0)
	}

	var mediafireLinkList []string
	doc.Find("a.shortc-button.medium.green").Each(func(i int, s *goquery.Selection) {
		mediafireLink, _ := s.Attr("href")
		mediafireLinkList = append(mediafireLinkList, mediafireLink)
		fmt.Println("    " + mediafireLink)
	})

	var imgLinkList []string
	imgLinkList = getDetailPageImgList(url)

	//这里创建一个需要写入的map
	dataMap := make(map[string]interface{})
	//将数据写入map
	dataMap["title"] = title
	dataMap["updated"] = doc.Find(".updated").Text()
	dataMap["url"] = url
	dataMap["mediafireLink"] = mediafireLinkList
	dataMap["imgLinkList"] = imgLinkList
	//打开文件
	outputFile, _ := os.OpenFile(jsonFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(outputFile)
	//创建encoder 数据输出到file中
	encoder := json.NewEncoder(outputFile)
	//把dataMap的数据encode到file中
	err = encoder.Encode(dataMap)
	//异常处理
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

func getDetailPageImgList(detailPage string) (imgList []string) {
	fmt.Println("    " + detailPage)
	doc, err := goquery.NewDocument(detailPage)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find("#fukie2 img.aligncenter").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("src")
		imgList = append(imgList, val)
		fmt.Println("      " + val)
	})

	nextDetailPage, _ := doc.Find("head link[rel=next]").Attr("href")
	if nextDetailPage != "" {
		temp := getDetailPageImgList(nextDetailPage)
		imgList = append(imgList, temp...)
	}

	return
}
