package xgmn01_test

import (
	"fmt"
	"page.github.io/pkg/xgmn01"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestReadyToUpload
func TestReadyToUpload(t *testing.T) {
	// step1
	firstUrl := xgmn01.Domain + "/Xgyw"
	jsonFiles := xgmn01.DownloadToJson(firstUrl)
	if len(jsonFiles) == 0 {
		return
	}

	// step2
	imgFiles := xgmn01.DownloadFromJson(jsonFiles)
	if len(imgFiles) == 0 {
		return
	}

	// step3
	xgmn01.ImgToVideo(imgFiles)
}

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestDownloadToJson
func TestDownloadToJson(t *testing.T) {
	firstUrl := xgmn01.Domain + "/Xgyw"
	xgmn01.DownloadToJson(firstUrl)
}

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestDownloadFromJson
func TestDownloadFromJson(t *testing.T) {
	var jsonFiles []string
	_ = xgmn01.DownloadFromJson(jsonFiles)
}

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestImgToVideo
func TestImgToVideo(t *testing.T) {
	var imgDirs []string
	xgmn01.ImgToVideo(imgDirs)
}
