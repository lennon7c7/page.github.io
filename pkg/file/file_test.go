package file_test

import (
	"fmt"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"path"
	"runtime"
	"testing"
	"time"
)

// go test -v pkg/file/file_test.go
func TestGetFileList(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	pathName := "../../images/test/1"
	files := file.GetFiles(pathName)
	for _, f := range files {
		fmt.Println(f)
	}

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/file/file_test.go -run TestSerialRename
func TestSerialRename(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(filename)))
	pathName := root + "/images"
	files := img.GetFiles(pathName)

	file.SerialRename(files)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/file/file_test.go -run TestGetRedirectUrl
func TestGetRedirectUrl(t *testing.T) {
	input := "https://www.xgmn02.com/uploadfile/202205/28/FF13279573.jpg"
	output := file.GetRedirectUrl(input)

	expect := "1https://j.20dh.top/Uploadfile/202205/28/FF13279573.jpg"
	if expect != output {
		t.FailNow()
	}
}
