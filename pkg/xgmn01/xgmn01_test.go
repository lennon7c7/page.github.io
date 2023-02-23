package xgmn01_test

import (
	"page.github.io/pkg/xgmn01"
	"runtime"
	"testing"
)

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestXgywImg
func TestXgywImg(t *testing.T) {
	runtime.GOMAXPROCS(4)

	firstUrl := xgmn01.Domain + "/Xgyw/"
	xgmn01.ListPage(firstUrl)
}
