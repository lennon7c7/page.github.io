package xgmn01_test

import (
	"fmt"
	"page.github.io/pkg/xgmn01"
	"runtime"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -timeout 0 -v pkg/xgmn01/xgmn01_test.go -run TestXgywImg
func TestXgywImg(t *testing.T) {
	runtime.GOMAXPROCS(4)

	firstUrl := xgmn01.Domain + "/Xgyw/page_118.html"
	xgmn01.ListPage(firstUrl)
}
