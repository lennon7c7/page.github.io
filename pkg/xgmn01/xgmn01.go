package xgmn01

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"page.github.io/pkg/file"
	"strings"
)

var Channel chan int
var Domain = "https://www.xgmn02.com"
var BaseDownloadJsonPath = "../../json/" + file.GetNameWithoutExt() + "/"

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
