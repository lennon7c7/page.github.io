package xgmn01

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Channel chan int
var Domain = "https://www.xgmn02.com"
var BaseDownloadJsonPath = "../../json/" + file.GetNameWithoutExt() + "/"
var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"

type jsonData struct {
	Title       string
	Updated     string
	Url         string
	ImgLinkList []string
}

func DownloadFromJson() {
	runtime.GOMAXPROCS(10)
	extName := ".json"
	err := filepath.Walk(BaseDownloadJsonPath, func(pathFile string, info os.FileInfo, err error) error {
		if extName != path.Ext(pathFile) {
			return nil
		}

		content, err := os.ReadFile(pathFile)
		if err != nil {
			fmt.Println(pathFile, err)
			return err
		}

		//Now let's unmarshall the data into `payload`
		var payload jsonData
		err = json.Unmarshal(content, &payload)
		if err != nil {
			fmt.Println(pathFile, err)
			return err
		}

		if len(payload.ImgLinkList) == 0 {
			return nil
		}

		downloadImgPath := BaseDownloadImgPath + payload.Updated + "/" + payload.Title
		dirEntries, _ := os.ReadDir(downloadImgPath)
		if file.Exists(downloadImgPath) && len(dirEntries) > 0 {
			return nil
		}

		err = os.MkdirAll(downloadImgPath, 0777)
		if err != nil {
			return err
		}

		fmt.Println("  " + payload.Url)
		for i, imgLink := range payload.ImgLinkList {
			filename := downloadImgPath + fmt.Sprintf("/%04d", i) + path.Ext(imgLink)
			Channel = make(chan int)
			go downloadImage(imgLink, filename)
			<-Channel
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}

func downloadImage(downloadLink string, filename string) {
	defer func() {
		Channel <- 0
	}()

	err := got.New().Download(downloadLink, filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = img.CutBorder(filename, filename, 100)

	output := strconv.FormatInt(time.Now().UnixNano(), 10) + path.Ext(filename)
	outputWidth := 720
	outputHeight := 100
	_ = img.Cut(filename, output, outputWidth, outputHeight)

	exists := img.IsWatermark(output)
	if exists {
		_ = os.Remove(filename)
	}
}

func ListPage(url string) {
	fmt.Println(url)
	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".widget-title a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		if detailUrl != "" {
			Channel = make(chan int)
			go detailPage(Domain + detailUrl)
			<-Channel
		}
	})

	nextUrl, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextUrl != "" {
		nextUrl = Domain + nextUrl
		ListPage(nextUrl)
	}
}

func detailPage(url string) {
	fmt.Println("  " + url)

	defer func() {
		Channel <- 0
	}()

	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find(".article-title").First().Text()

	updated := doc.Find(".article-meta .item-1").First().Text()
	updated = strings.Replace(updated, "更新：", "", 1)
	updated = strings.Replace(updated, ".", "/", 2)

	jsonFile := BaseDownloadJsonPath + updated + "/" + title + ".json"

	var imgLinkList []string
	imgLinkList = getDetailPageImgList(url)

	dataMap := make(map[string]interface{})
	dataMap["title"] = title
	dataMap["updated"] = updated
	dataMap["url"] = url
	dataMap["imgLinkList"] = imgLinkList

	err = file.Create(jsonFile, dataMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func getDetailPageImgList(detailPage string) (imgList []string) {
	//fmt.Println("    " + detailPage)

	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(detailPage)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".article-content img").Each(func(i int, s *goquery.Selection) {
		imgSrc, _ := s.Attr("src")
		imgSrc = Domain + imgSrc
		imgSrc = file.GetRedirectUrl(imgSrc)

		imgList = append(imgList, imgSrc)
		//fmt.Println("      " + imgSrc)
	})

	nextPage, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextPage != "" {
		nextPage = Domain + nextPage
		temp := getDetailPageImgList(nextPage)
		imgList = append(imgList, temp...)
	}

	return
}
