package xgmn01

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/melbahja/got"
	"os"
	"page.github.io/pkg/ffmpeg"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var Channel chan int
var ChannelString chan string
var IsExit bool
var Domain = "https://www.xgmn02.com"
var BaseDownloadJsonPath = "../../json/" + file.GetNameWithoutExt() + "/"
var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"
var BaseOutputVideoPath = "../../video/" + file.GetNameWithoutExt() + "/"

type jsonData struct {
	Title       string
	Updated     string
	Url         string
	ImgLinkList []string
}

func DownloadToJson(url string) (files []string) {
	runtime.GOMAXPROCS(5)

	IsExit = false
	fmt.Println(url)
	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".widget-title a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		if detailUrl == "" {
			return
		}

		ChannelString = make(chan string)
		go detailPage(Domain + detailUrl)
		tempFile := <-ChannelString
		if tempFile != "" {
			files = append(files, tempFile)
		}
	})

	if IsExit {
		return
	}

	nextUrl, _ := doc.Find(".pagination a:contains(下一页)").First().Attr("href")
	if nextUrl != "" {
		nextUrl = Domain + nextUrl
		temp := DownloadToJson(nextUrl)
		files = append(files, temp...)
	}

	return
}

func detailPage(url string) {
	jsonFile := ""
	fmt.Println("  ", url)

	defer func() {
		ChannelString <- jsonFile
	}()

	//goland:noinspection GoDeprecation
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find(".article-title").First().Text()
	if title == "" {
		return
	}

	updated := doc.Find(".article-meta .item-1").First().Text()
	updated = strings.Replace(updated, "更新：", "", 1)
	updated = strings.Replace(updated, ".", "/", 2)
	if updated == "" {
		return
	}

	jsonFile = BaseDownloadJsonPath + updated + "/" + title + ".json"
	if file.Exists(jsonFile) {
		fmt.Println("---------- no shit ---------- ")
		IsExit = true
		jsonFile = ""
		return
	}

	var imgLinkList []string
	imgLinkList = getDetailPageImgList(url)
	if len(imgLinkList) == 0 {
		jsonFile = ""
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["title"] = title
	dataMap["updated"] = updated
	dataMap["url"] = url
	dataMap["imgLinkList"] = imgLinkList

	err = file.Create(jsonFile, dataMap)
	if err != nil {
		fmt.Println(err)
		jsonFile = ""
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

func DownloadFromJson(jsonFiles []string) (imgDirs []string) {
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

	runtime.GOMAXPROCS(10)
	for _, jsonFile := range jsonFiles {
		imgDir, _ := downloadFromJson(jsonFile)
		if imgDir != "" {
			imgDirs = append(imgDirs, imgDir)
		}
	}

	return
}

func downloadFromJson(jsonFile string) (downloadImgPath string, err error) {
	content, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Now let's unmarshall the data into `payload`
	var payload jsonData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(payload.ImgLinkList) == 0 {
		fmt.Println("len(payload.ImgLinkList) == 0")
		return
	}

	downloadImgPath = BaseDownloadImgPath + payload.Updated + "/" + payload.Title
	dirEntries, _ := os.ReadDir(downloadImgPath)
	if file.Exists(downloadImgPath) && len(dirEntries) > 0 {
		fmt.Println("---------- no shit ---------- ")
		IsExit = true
		return
	}

	err = os.MkdirAll(downloadImgPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("  ", payload.Url)
	for i, imgLink := range payload.ImgLinkList {
		filename := downloadImgPath + fmt.Sprintf("/%04d", i) + path.Ext(imgLink)
		Channel = make(chan int)
		go downloadImage(imgLink, filename)
		<-Channel
	}
	file.SerialRename(img.GetFiles(downloadImgPath))
	img.BatchMaxImageWidthHeight(downloadImgPath)

	return
}

func ImgToVideo(imgDirs []string) {
	if len(imgDirs) == 0 {
		infos, err := os.ReadDir(BaseDownloadImgPath)
		if err != nil {
			return
		}

		for _, yearDir := range infos {
			if !yearDir.IsDir() {
				continue
			}

			yearDirAbs, _ := filepath.Abs(BaseDownloadImgPath + "/" + yearDir.Name())
			monthDirEntries, _ := os.ReadDir(yearDirAbs)
			for _, monthDir := range monthDirEntries {
				if !monthDir.IsDir() {
					continue
				}

				dayDirAbs, _ := filepath.Abs(yearDirAbs + "/" + monthDir.Name())
				dayDirEntries, _ := os.ReadDir(dayDirAbs)
				for _, dayDir := range dayDirEntries {
					if !dayDir.IsDir() {
						continue
					}

					blogDirAbs, _ := filepath.Abs(dayDirAbs + "/" + dayDir.Name())
					blogDirEntries, _ := os.ReadDir(blogDirAbs)
					for _, blogDir := range blogDirEntries {
						if !blogDir.IsDir() {
							continue
						}

						inputImgDir := blogDirAbs + "/" + blogDir.Name()
						inputImgList := img.GetFiles(inputImgDir)
						if len(inputImgList) <= 3 {
							continue
						}

						output := BaseOutputVideoPath + yearDir.Name() + "/" + monthDir.Name() + "/" + dayDir.Name() + "/" + blogDir.Name() + ".mp4"
						if file.Exists(output) {
							continue
						}

						imgDirs = append(imgDirs, inputImgDir)
					}
				}
			}
		}

		// 降序
		sort.Sort(sort.Reverse(sort.StringSlice(imgDirs)))
	}

	for _, inputImgDir := range imgDirs {
		blogDirName := filepath.Base(inputImgDir)
		dayDirName := filepath.Base(filepath.Dir(inputImgDir))
		monthDirName := filepath.Base(filepath.Dir(filepath.Dir(inputImgDir)))
		yearDirName := filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(inputImgDir))))

		inputImgTemplate := inputImgDir + "/%04d.jpg"
		inputImgList := img.GetFiles(inputImgDir)
		if len(inputImgList) <= 3 {
			continue
		}

		output := BaseOutputVideoPath + yearDirName + "/" + monthDirName + "/" + dayDirName + "/" + blogDirName + ".mp4"
		if file.Exists(output) {
			fmt.Println("---------- no shit ---------- ")
			continue
		}
		ffmpeg.Img2Video(inputImgTemplate, output)

		inputVideo := output
		inputAudio := file.GetRandomAudio()
		output = inputVideo
		ffmpeg.AddAudio2Video(inputVideo, inputAudio, output)
		fmt.Println(filepath.Abs(output))
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
