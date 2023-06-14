package img_test

import (
	"fmt"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -v pkg/img/img_test.go -run TestMaxImageWidthHeight
func TestMaxImageWidthHeight(t *testing.T) {
	//canvasWidth := 2048
	//canvasHeight := 2000
	//imgFile := "../../images/nature-1.jpg"
	//img.MaxImageWidthHeight(canvasWidth, canvasHeight, imgFile)

	pathName := "../../images/test"
	files := img.GetFiles(pathName)

	x, y := img.GetMaxWidthHeight(files)

	for _, tempFile := range files {
		img.MaxImageWidthHeight(x, y, tempFile)
	}
}

// go test -v pkg/img/img_test.go -run TestGetMaxWidthHeight
func TestGetMaxWidthHeight(t *testing.T) {
	pathName := "../../images/test/1"
	files := file.GetFiles(pathName)

	x, y := img.GetMaxWidthHeight(files)
	fmt.Println(x, y)
}

// go test -v pkg/img/img_test.go -run TestResize
func TestResize(t *testing.T) {
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
}

// go test -v pkg/img/img_test.go -run TestCut
func TestCut(t *testing.T) {
	input := "../../images/test/i-am-watermark.jpg"
	output := "../../images/test/watermark-after-clear.jpg"
	outputWidth := 720
	outputHeight := 100
	_ = img.Cut(input, output, outputWidth, outputHeight)

	//input = "../../images/test/0031.jpg"
	//output = "../../images/test/0031-.jpg"
	//outputWidth = 720
	//outputHeight = 100
	//img.Cut(input, output, outputWidth, outputHeight)
	//
	//input = "../../images/test/0036.jpg"
	//output = "../../images/test/0036-.jpg"
	//outputWidth = 720
	//outputHeight = 100
	//img.Cut(input, output, outputWidth, outputHeight)
}

// go test -v pkg/img/img_test.go -run TestRename
func TestRename(t *testing.T) {
	pathName := "../../images/test/1"
	files := img.GetFiles(pathName)
	file.SerialRename(files)
}

// go test -v pkg/img/img_test.go -run TestIsWatermark
func TestIsWatermark(t *testing.T) {
	input := "../../images/test/i-am-watermark.jpg"
	input = "../../images/test/1/0000.jpeg"
	output := "../../images/test/is-watermark.jpg"
	outputWidth := 720
	outputHeight := 100
	_ = img.Cut(input, output, outputWidth, outputHeight)

	fmt.Println(img.IsWatermark(output))
}

// go test -v pkg/img/img_test.go -run TestCutBorder
func TestCutBorder(t *testing.T) {
	input := "../../images/test/i-am-watermark.jpg"
	input = "../../images/test/1/0000.jpeg"
	input = "../../images/test/0000-0000.jpg"
	output := "../../images/test/is-watermark.jpg"
	output = input
	border := 100
	_ = img.CutBorder(input, output, border)
}

func TestUrl2File(t *testing.T) {
	input := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"
	output := "../../images/baidu.png"
	err := img.Url2File(input, output)
	fmt.Println(err)
}

func TestGenerateMask(t *testing.T) {
	imgWidth := 1024
	imgHeight := 1024
	// bbox: [100, 200, 597, 104]
	maskX := 332
	maskY := 136
	maskWidth := 359
	maskHeight := 499
	t.Log(img.GenerateRectMask(imgWidth, imgHeight, maskX, maskY, maskWidth, maskHeight))
}
