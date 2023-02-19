package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"page.github.io/pkg/file"
	"path"
	"runtime"
	"strings"
	"time"
)

// 一个妹子图片网站
var domain = "https://www.xgmn01.com"
var baseDownloadPath = "./images/" + file.GetNameWithoutExt() + "/"

var c chan int

func main() {
	runtime.GOMAXPROCS(4)

	var url = domain + "/new.html"
	listPage(url)
}

// url->Document->所有图片url->开启多线程进行下载->保存到本地
func listPage(url string) {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".widget-title a").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Attr("title")

		downloadPath := s.Parent().Find(".post-like").Text()
		downloadPath = strings.Replace(downloadPath, "-", "/", 2)
		downloadPath += "/" + title
		downloadPath = baseDownloadPath + downloadPath + "/"

		if file.Exists(downloadPath) {
			return
		}

		url, _ = s.Attr("href")
		if url != "" {
			detailPage(url, 0)
		}
	})
}

func detailPage(url string, page int) {
	fmt.Println("  " + url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find(".article-title").First().Text()

	downloadPath := doc.Find(".article-meta .item-1").First().Text()
	downloadPath = strings.Replace(downloadPath, "更新：", "", 1)
	downloadPath = strings.Replace(downloadPath, ".", "/", 2)
	downloadPath += "/" + title
	downloadPath = baseDownloadPath + downloadPath + "/"

	doc.Find(".article-content img").Each(func(i int, s *goquery.Selection) {
		url, _ = s.Attr("src")
		url = strings.Replace(url, "uploadfile", "Uploadfile", 1)

		pageString := fmt.Sprintf("%04d-", page)
		filename := downloadPath + pageString + fmt.Sprintf("%04d", i) + path.Ext(url)

		c = make(chan int)
		go downloadImage(domain+url, filename)
		<-c
	})

	nextUrl, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextUrl != "" {
		nextUrl = domain + nextUrl
		nextPage := page
		nextPage++
		detailPage(nextUrl, nextPage)
	}
}

// 根据url 创建http 请求的 request
// 网站有反爬虫策略 wireshark 不解释
func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return req
	}

	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36")
	//req.Header.Set("Referer", domain)
	return req
}

// 下载图片
func downloadImage(url string, filename string) {
	fmt.Println("    " + url)

	req := buildRequest(url)
	http.DefaultClient.Timeout = 10 * time.Second
	resp, err := http.DefaultClient.Do(req)
	defer func(err error) {
		c <- 0
		if err != nil {
			return
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}(err)

	if err != nil {
		fmt.Println(err)
		fmt.Println("      failed download")
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("      resp.StatusCode: " + string(rune(resp.StatusCode)))
		fmt.Println("      http.StatusOK: " + string(rune(http.StatusOK)))
		fmt.Println("      failed download")
		return
	}

	downloadPath := path.Dir(filename)
	err = os.MkdirAll(downloadPath, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	localFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	if _, err := io.Copy(localFile, resp.Body); err != nil {
		fmt.Println(err)
		fmt.Println("      failed save " + url)
		return
	}
}
