package mrcong_test

import (
	"fmt"
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
