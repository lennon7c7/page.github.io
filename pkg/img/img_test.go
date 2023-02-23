package img_test

import (
	"fmt"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"testing"
	"time"
)

// go test -v pkg/img/img_test.go -run TestMaxImageWidthHeight
func TestMaxImageWidthHeight(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	//canvasWidth := 2048
	//canvasHeight := 2000
	//imgFile := "../../images/nature-1.jpg"
	//img.MaxImageWidthHeight(canvasWidth, canvasHeight, imgFile)

	pathName := "../../images/test/1"
	files := file.GetFiles(pathName)

	x, y := img.GetMaxWidthHeight(files)

	for _, tempFile := range files {
		img.MaxImageWidthHeight(x, y, tempFile)
	}

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/img/img_test.go -run TestGetMaxWidthHeight
func TestGetMaxWidthHeight(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	pathName := "../../images/test/1"
	files := file.GetFiles(pathName)

	x, y := img.GetMaxWidthHeight(files)
	fmt.Println(x, y)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/img/img_test.go -run TestResize
func TestResize(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	//input := "../../images/test/1/Coser-No.085-MrCong.com-005.jpg"
	//output := "../../images/test/1/Coser-No.085-MrCong.com-005--.jpg"
	//output = input
	//img.AuthThumbnail(input, output)
	//
	//input = "../../images/test/1/Coser-No.085-MrCong.com-006.jpg"
	//output = "../../images/test/1/Coser-No.085-MrCong.com-006--.jpg"
	//output = input
	//img.AuthThumbnail(input, output)

	pathName := "../../images/test/1"
	files := img.GetFiles(pathName)
	for _, tempFile := range files {
		img.AuthThumbnail(tempFile, tempFile)
	}

	x, y := img.GetMaxWidthHeight(files)
	for _, tempFile := range files {
		img.MaxImageWidthHeight(x, y, tempFile)
	}

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/img/img_test.go -run TestCut
func TestCut(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	input := "../../images/test/0013.jpg"
	output := "../../images/test/0013-.jpg"
	outputWidth := 720
	outputHeight := 90
	img.Cut(input, output, outputWidth, outputHeight)

	input = "../../images/test/0031.jpg"
	output = "../../images/test/0031-.jpg"
	outputWidth = 720
	outputHeight = 90
	img.Cut(input, output, outputWidth, outputHeight)

	input = "../../images/test/0036.jpg"
	output = "../../images/test/0036-.jpg"
	outputWidth = 720
	outputHeight = 90
	img.Cut(input, output, outputWidth, outputHeight)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/img/img_test.go -run TestRename
func TestRename(t *testing.T) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	pathName := "../../images/test/1"
	files := img.GetFiles(pathName)
	file.SerialRename(files)

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}
