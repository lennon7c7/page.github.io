package mrcong_test

import (
	"fmt"
	"page.github.io/pkg/ffmpeg"
	"page.github.io/pkg/file"
	"page.github.io/pkg/mrcong"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -timeout 0 -v pkg/mrcong/mrcong_test.go -run TestDownloadToJson
func TestDownloadToJson(t *testing.T) {
	webUrl := mrcong.Domain + "/tag/%e5%b0%8f%e4%bb%93%e5%8d%83%e4%bb%a3w/"
	mrcong.DownloadToJson(webUrl)
}

// go test -timeout 0 -v pkg/mrcong/mrcong_test.go -run TestDownloadMediafireLink
func TestDownloadMediafireLink(t *testing.T) {
	var jsonFiles []string
	mrcong.DownloadMediafireLink(jsonFiles)
}

// go test -timeout 0 -v pkg/mrcong/mrcong_test.go -run TestDownloadByTag
func TestDownloadByTag(t *testing.T) {
	// step1
	webUrl := mrcong.Domain + "/tag/%e7%99%bd%e9%93%b681/page/4/"
	//webUrl = mrcong.Domain + "/tag/merry/"
	jsonFiles := mrcong.ListPage(webUrl)
	if len(jsonFiles) == 0 {
		return
	}

	// step2
	zipFiles := mrcong.DownloadMediafireLink(jsonFiles)
	if len(zipFiles) == 0 {
		return
	}

	// step3
	pathName := "../../zip/mrcong"
	files := file.GetFiles(pathName)
	outpuf := "../../images/test"
	for _, f := range files {
		mrcong.Unrar(f, outpuf)
	}

	// step4
	files = file.GetDirs(outpuf)
	for _, f := range files {
		outputVideo := ""
		t.Log(ffmpeg.ConcatVideo2Video(f, outputVideo))
	}
}
