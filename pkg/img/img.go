package img

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"page.github.io/pkg/array"
	"page.github.io/pkg/file"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var ExtList = []string{".jpg", ".jpeg", ".png", ".webp"}

type ApiWatermarkData struct {
	Base64 string `json:"base64"`
	Result []struct {
		SavePath string `json:"save_path"`
		Data     []struct {
			Text            string  `json:"text"`
			Confidence      float64 `json:"confidence"`
			TextBoxPosition [][]int `json:"text_box_position"`
		} `json:"data"`
	} `json:"result"`
}

// MaxImageWidthHeight 扩张图片宽高
// 扩张之处以黑色背景填充
func MaxImageWidthHeight(canvasWidth int, canvasHeight int, imgFile string) {
	// 返回一个矩形
	rectangle := image.Rect(0, 0, canvasWidth, canvasHeight)
	rgba := image.NewRGBA(rectangle)

	// 创建一个新的上下文
	context := freetype.NewContext()
	context.SetDPI(70)                                                         // 设置屏幕分辨率，单位为每英寸点数。
	context.SetClip(rgba.Bounds())                                             //设置用于绘制的剪辑矩形。
	context.SetDst(rgba)                                                       //设置绘制操作的目标图像。
	context.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 1})) //设置用于绘制操作的源图像

	// 图片水印
	img, _ := os.Open(imgFile)
	defer func(img *os.File) {
		err := img.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(img)
	img1, _ := jpeg.Decode(img) // 读取一个JPEG图像并将其作为image.Image返回

	offsetX := canvasWidth/2 - img1.Bounds().Dx()/2
	offsetY := canvasHeight/2 - img1.Bounds().Dy()/2
	offset := image.Pt(offsetX, offsetY)
	if canvasWidth == img1.Bounds().Dx() {
		offset = image.Pt(0, offsetY)
	} else if canvasHeight == img1.Bounds().Dy() {
		offset = image.Pt(offsetX, 0)
	}
	draw.Draw(rgba, img1.Bounds().Add(offset), img1, image.Point{}, draw.Over)

	// 创建图片
	tempFile, err := os.Create(imgFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 将图像写入file
	err = jpeg.Encode(tempFile, rgba, &jpeg.Options{Quality: 100})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile)
}

func GetMaxWidthHeight(files []string) (width int, height int) {
	for _, tempFile := range files {
		img, _ := os.Open(tempFile)
		img1, _ := jpeg.Decode(img)

		if img1.Bounds().Dx() > width {
			width = img1.Bounds().Dx()
		}

		if img1.Bounds().Dy() > height {
			height = img1.Bounds().Dy()
		}
	}

	return
}

func GetMaxFilename(pathName string) (filenameInt int) {
	filenameMax := 0
	files := GetFiles(pathName)
	for _, tempFile := range files {
		// 获取文件名带后缀
		filenameWithSuffix := filepath.Base(tempFile)
		// 获取文件后缀
		fileSuffix := path.Ext(filenameWithSuffix)
		// 获取文件名
		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)

		// 把文件名转换成数字
		filenameInt, _ = strconv.Atoi(filenameOnly)
		if filenameInt > filenameMax {
			filenameMax = filenameInt
		}
	}

	return
}

func GetFiles(pathName string) (files []string) {
	if !file.IsDir(pathName) {
		pathName = path.Dir(pathName)
	}

	tempFiles := file.GetFiles(pathName)

	for _, value := range tempFiles {
		exists, _ := array.IsExists(path.Ext(value), ExtList)
		if !exists {
			continue
		}

		files = append(files, value)
	}

	return
}

func AuthThumbnail(inputFile string, outputFile string) {
	tempFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// decode jpeg into image.Image
	img1, err := jpeg.Decode(tempFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tempFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	maxWidth := 1920
	maxHeight := 1920
	m := resize.Thumbnail(uint(maxWidth), uint(maxHeight), img1, resize.Lanczos2)

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(out)

	// write new image to tempFile
	err = jpeg.Encode(out, m, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Cut(inputFile string, outputFile string, outputWidth int, outputHeight int) (err error) {
	rectangle := image.Rect(0, 0, outputWidth, outputHeight)
	rgba := image.NewRGBA(rectangle)

	img, _ := os.Open(inputFile)
	defer func(img *os.File) {
		err := img.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(img)
	img1, err := jpeg.Decode(img)
	if err != nil {
		return
	}

	offsetX := outputWidth/2 - img1.Bounds().Dx()/2
	offsetY := outputHeight/2 - img1.Bounds().Dy()/2
	offset := image.Pt(offsetX, offsetY)
	if outputWidth == img1.Bounds().Dx() {
		offset = image.Pt(0, offsetY)
	} else if outputHeight == img1.Bounds().Dy() {
		offset = image.Pt(offsetX, 0)
	}

	draw.Draw(rgba, img1.Bounds().Add(offset), img1, image.Point{}, draw.Over)

	tempFile, err := os.Create(outputFile)
	if err != nil {
		return
	}

	err = jpeg.Encode(tempFile, rgba, &jpeg.Options{Quality: 70})
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile)

	return
}

func CutBorder(inputFile string, outputFile string, border int) (err error) {
	img, _ := os.Open(inputFile)
	defer func(img *os.File) {
		err := img.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(img)
	img1, err := jpeg.Decode(img)
	if err != nil {
		return
	}

	outputWidth := img1.Bounds().Dx() - border
	if outputWidth%2 != 0 {
		outputWidth--
	}
	outputHeight := img1.Bounds().Dy() - border
	if outputHeight%2 != 0 {
		outputHeight--
	}
	err = Cut(inputFile, outputFile, outputWidth, outputHeight)

	return
}

func IsWatermark(inputFile string) (exists bool) {
	watermarks := []string{"www.", ".net"}
	exists = false

	srcByte, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	base64String := base64.StdEncoding.EncodeToString(srcByte)

	url := "https://www.paddlepaddle.org.cn/paddlehub-api/image_classification/chinese_ocr_db_crnn_mobile"
	method := "POST"
	payload := strings.NewReader(`{"image":"` + base64String + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	_ = os.Remove(inputFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var resultData ApiWatermarkData
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(resultData.Result) == 0 || len(resultData.Result[0].Data) == 0 {
		return
	}

	text := resultData.Result[0].Data[0].Text
	text = strings.ToLower(text)
	for _, v := range watermarks {
		if strings.Contains(text, v) {
			return true
		}
	}

	return
}

func BatchMaxImageWidthHeight(dirName string) {
	files := GetFiles(dirName)
	x, y := GetMaxWidthHeight(files)
	for _, tempFile := range files {
		MaxImageWidthHeight(x, y, tempFile)
	}

}
