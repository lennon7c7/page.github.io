package mrcong_test

import (
	"page.github.io/pkg/mrcong"
	"testing"
)

// go test -timeout 0 -v pkg/mrcong/mrcong_test.go -run TestDownloadToJson
func TestDownloadToJson(t *testing.T) {
	mrcong.DownloadToJson()
}
