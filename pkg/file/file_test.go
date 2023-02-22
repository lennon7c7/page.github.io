package file_test

import (
	"fmt"
	"page.github.io/pkg/file"
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

	pathName := "../../images/test/1"
	files := file.GetFiles(pathName)

	file.SerialRename(files)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}
