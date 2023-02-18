package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var baseDownloadJsonPath = "./json/" + getCurrentRuntimeFilename() + "/"

func main() {
	runtime.GOMAXPROCS(4)

	// 一个妹子图片网站
	var listUrl = "https://mrcong.com"
	mrcongListPage(listUrl)
}

func mrcongListPage(url string) {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".post-listing .post-box-title a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		if detailUrl != "" {
			mrcongDetailPage(detailUrl, 0)
		}
	})

	nextUrl, _ := doc.Find("head link[rel=next]").Attr("href")
	if nextUrl != "" {
		mrcongListPage(nextUrl)
	}
}

func mrcongDetailPage(url string, page int) {
	fmt.Println("  " + url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	title := doc.Find("#crumbs .current").Text()

	downloadPath := doc.Find(".updated").Text()
	downloadPath = strings.Replace(downloadPath, "-", "/", 2)
	downloadPath = baseDownloadJsonPath + downloadPath + "/"

	if !fileExists(downloadPath) {
		err = os.MkdirAll(downloadPath, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	jsonFile := downloadPath + title + ".json"
	if fileExists(jsonFile) {
		fmt.Println("---------- no shit ---------- ")
		os.Exit(0)
	}

	var mediafireLinkList []string
	doc.Find("a.shortc-button.medium.green").Each(func(i int, s *goquery.Selection) {
		mediafireLink, _ := s.Attr("href")
		mediafireLinkList = append(mediafireLinkList, mediafireLink)
		fmt.Println("    " + mediafireLink)
	})

	//这里创建一个需要写入的map
	dataMap := make(map[string]interface{})
	//将数据写入map
	dataMap["title"] = title
	dataMap["url"] = url
	dataMap["mediafireLink"] = mediafireLinkList
	//打开文件
	file, _ := os.OpenFile(jsonFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer file.Close()
	//创建encoder 数据输出到file中
	encoder := json.NewEncoder(file)
	//把dataMap的数据encode到file中
	err = encoder.Encode(dataMap)
	//异常处理
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getImgListFrom(inputPath string) (files []string) {
	names := []string{".jpg", ".png", ".jpeg"}
	err := filepath.Walk(inputPath, func(pathFile string, info os.FileInfo, err error) error {
		if inputPath == pathFile {
			return nil
		}

		exists, _ := inArray(path.Ext(pathFile), names)
		if !exists {
			return nil
		}

		files = append(files, pathFile)
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// get current runtime filename without ext
func getCurrentRuntimeFilename() string {
	_, fullFilename, _, _ := runtime.Caller(0)
	//获取文件名带后缀
	filenameWithSuffix := path.Base(fullFilename)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

	return filenameOnly
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func unrar(input string, output string) {
	command := "unrar x -pmrcong.com -inul -y " + input + " " + output
	_, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
}
