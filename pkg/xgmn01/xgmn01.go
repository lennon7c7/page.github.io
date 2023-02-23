package xgmn01

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"path"
	"strings"
)

var Channel chan int
var Domain = "https://xgmn01.com"
var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"

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
			downloadPath := detailPage(Domain+detailUrl, 0)
			file.SerialRename(img.GetFiles(downloadPath))
		}
	})

	nextUrl, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextUrl != "" {
		nextUrl = Domain + nextUrl
		ListPage(nextUrl)
	}
}

func detailPage(url string, page int) (downloadPath string) {
	fmt.Println("  " + url)
	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find(".article-title").First().Text()

	downloadPath = doc.Find(".article-meta .item-1").First().Text()
	downloadPath = strings.Replace(downloadPath, "更新：", "", 1)
	downloadPath = strings.Replace(downloadPath, ".", "/", 2)
	downloadPath += "/" + title
	downloadPath = BaseDownloadImgPath + downloadPath + "/"

	if file.Exists(downloadPath) && page == 0 {
		fmt.Println("---------- no shit ---------- ")
		os.Exit(0)
	} else {
		err = os.MkdirAll(downloadPath, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	doc.Find(".article-content img").Each(func(i int, s *goquery.Selection) {
		url, _ = s.Attr("src")
		url = strings.Replace(url, "uploadfile", "Uploadfile", 1)

		pageString := fmt.Sprintf("%04d-", page)
		filename := downloadPath + pageString + fmt.Sprintf("%04d", i) + path.Ext(url)

		Channel = make(chan int)
		go downloadImage(Domain+url, filename)
		<-Channel
	})

	nextUrl, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextUrl != "" {
		nextUrl = Domain + nextUrl
		nextPage := page
		nextPage++
		detailPage(nextUrl, nextPage)
	}

	return
}

func downloadImage(downloadLink string, filename string) {
	fmt.Println("    " + downloadLink)

	err := got.New().Download(downloadLink, filename)
	if err != nil {
		fmt.Println(err)
	}

	Channel <- 0
}
