package img

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/melbahja/got"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"page.github.io/pkg/array"
	"page.github.io/pkg/file"
	"page.github.io/pkg/log"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	// 图片水印
	img, _ := os.Open(imgFile)
	defer func(img *os.File) {
		err := img.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(img)
	img1, _ := jpeg.Decode(img) // 读取一个JPEG图像并将其作为image.Image返回

	if img1.Bounds().Dx() == canvasWidth && img1.Bounds().Dy() == canvasHeight {
		// 图片宽高相等，不需要扩张
		return
	}

	// 返回一个矩形
	rectangle := image.Rect(0, 0, canvasWidth, canvasHeight)
	rgba := image.NewRGBA(rectangle)

	// 创建一个新的上下文
	context := freetype.NewContext()
	context.SetDPI(70)                                                         // 设置屏幕分辨率，单位为每英寸点数。
	context.SetClip(rgba.Bounds())                                             //设置用于绘制的剪辑矩形。
	context.SetDst(rgba)                                                       //设置绘制操作的目标图像。
	context.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 1})) //设置用于绘制操作的源图像

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
		img, err := os.Open(tempFile)
		if err != nil {
			fmt.Println(tempFile, err)
			return
		}
		img1, err := jpeg.Decode(img)
		if err != nil {
			fmt.Println(tempFile, err)
			return
		}

		if img1.Bounds().Dx() > width {
			width = img1.Bounds().Dx()
		}

		if img1.Bounds().Dy() > height {
			height = img1.Bounds().Dy()
		}
	}

	// 图片宽、高必须是偶数
	if width%2 != 0 {
		width--
	}
	if height%2 != 0 {
		height--
	}

	return
}

func GetMaxFilename(pathName string) (filenameInt int) {
	filenameMax := 0

	if !file.Exists(pathName) {
		return
	}

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

	maxWidth = 4096
	maxHeight = 3112
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

func CutByBbox(originImgInterface image.Image, x int, y int, width int, height int) (cutImgInterface image.Image, base64Img string, err error) {
	// 裁剪出一个长方形
	rect := image.Rect(x, y, x+width, y+height)
	cutImgInterface = originImgInterface.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	//// 将裁剪后的图片保存为 JPG 文件
	//fileStruct, err := os.Create(outputFilename)
	//if err != nil {
	//	return
	//}
	//defer func() {
	//	closeErr := fileStruct.Close()
	//	if closeErr != nil {
	//		err = closeErr
	//	}
	//}()
	//err = jpeg.Encode(fileStruct, cutImgInterface, &jpeg.Options{Quality: 100})
	//if err != nil {
	//	return
	//}

	// Convert image to Base64 string
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, cutImgInterface)
	if err != nil {
		return
	}
	imgBytes := buffer.Bytes()
	base64Img = base64.StdEncoding.EncodeToString(imgBytes)

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
	if len(files) <= 1 {
		return
	}

	x, y := GetMaxWidthHeight(files)
	for _, tempFile := range files {
		MaxImageWidthHeight(x, y, tempFile)
	}

}

func Url2File(inputImgUrl string, outputImgFile string) (err error) {
	filePath := path.Dir(outputImgFile)
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return
	}

	err = got.New().Download(inputImgUrl, outputImgFile)
	if err != nil {
		return
	}

	return
}

func File2Base64(inputImgUrl string) (outputImgBase64 string, err error) {
	srcByte, err := os.ReadFile(inputImgUrl)
	if err != nil {
		return
	}

	outputImgBase64 = base64.StdEncoding.EncodeToString(srcByte)

	return
}

func Http2Base64(inputImgUrl string) (outputImgBase64 string, err error) {
	// 发送 HTTP GET 请求，获取网络图片
	resp, err := http.Get(inputImgUrl)
	if err != nil {
		return
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	// 读取图片内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 将图片内容进行 Base64 编码
	outputImgBase64 = base64.StdEncoding.EncodeToString(data)

	return
}

func GenerateRectMask(imgWidth int, imgHeight int, maskX int, maskY int, maskWidth int, maskHeight int) (outputImgBase64 string, err error) {
	// 创建一个黑色背景的图像
	rect := image.Rect(0, 0, imgWidth, imgHeight)

	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	// 在图像上绘制长方形
	draw.Draw(img, image.Rect(maskX, maskY, maskX+maskWidth, maskY+maskHeight), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// 将图像保存为PNG文件
	output := fmt.Sprintf("./%d.png", time.Now().Unix())
	filePoint, err := os.Create(output)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}

		_ = os.Remove(output)
	}(filePoint)

	// 将图像转换为字节切片
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return
	}

	// 将字节切片转换为 Base64 编码的字符串
	outputImgBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	err = png.Encode(filePoint, img)
	if err != nil {
		return
	}

	return
}

// GetImageSizeFromBase64 获取 Base64 编码图片的宽度和高度
func GetImageSizeFromBase64(base64String string) (width int, height int, err error) {
	substr := ","
	if strings.Contains(base64String, substr) {
		// 兼容
		base64String = strings.Split(base64String, substr)[1]
	}

	// 将 Base64 编码转换为字节数组
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	// 获取图片宽度和高度
	width = img.Bounds().Dx()
	height = img.Bounds().Dy()

	return
}

// Base64ToFile 将 Base64 编码的字符串转换为图片文件
func Base64ToFile(base64String string, filePath string) (err error) {
	substr := ","
	if strings.Contains(base64String, substr) {
		// 兼容
		base64String = strings.Split(base64String, substr)[1]
	}

	// Decode Base64 string to byte slice
	imgData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// Create new fileStruct
	fileStruct, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := fileStruct.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	// Write byte slice to fileStruct
	_, err = fileStruct.Write(imgData)
	if err != nil {
		return err
	}

	// Confirm that image fileStruct was created successfully
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(fileData) != len(imgData) {
		return errors.New("image fileStruct was not created correctly")
	}

	return nil
}

// Base64ToImgInterface 将 Base64 编码的字符串转换为图片对象
func Base64ToImgInterface(base64String string) (imgInterface image.Image, err error) {
	substr := ","
	if strings.Contains(base64String, substr) {
		// 兼容
		base64String = strings.Split(base64String, substr)[1]
	}

	// Decode Base64 string
	imgBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return
	}

	// Check the first few bytes of the decoded data to determine the image format
	if len(imgBytes) < 4 {
		return nil, io.EOF
	}

	// Decode image from bytes
	if imgBytes[0] == 0xff && imgBytes[1] == 0xd8 && imgBytes[2] == 0xff {
		imgInterface, err = jpeg.Decode(strings.NewReader(string(imgBytes)))
		if err != nil {
			return
		}
	} else if imgBytes[0] == 0x89 && imgBytes[1] == 0x50 && imgBytes[2] == 0x4e && imgBytes[3] == 0x47 {
		imgInterface, err = png.Decode(strings.NewReader(string(imgBytes)))
		if err != nil {
			return
		}
	} else {
		return nil, errors.New("unknown image type")
	}

	return
}

// DrawTransparentBackground 将裁剪后的图片放入一个新的图片中
func DrawTransparentBackground(cutImgBase64 string, backgroundWidth int, backgroundHeight int, maskX int, maskY int, outputFilename string) (err error) {
	cutImgInterface, err := Base64ToImgInterface(cutImgBase64)
	if err != nil {
		return
	}

	// 创建一个 w * h 的 RGBA 图片
	rect := image.Rect(0, 0, backgroundWidth, backgroundHeight)
	newImg := image.NewRGBA(rect)

	// 将新图片中所有像素设置为黑色
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	// 在新图片中绘制原图
	draw.Draw(newImg, cutImgInterface.Bounds().Add(image.Pt(maskX, maskY)), cutImgInterface, cutImgInterface.Bounds().Min, draw.Over)

	// 将图片保存为 JPG 文件
	out, err := os.Create(outputFilename)
	if err != nil {
		return
	}
	defer func() {
		closeErr := out.Close()
		if closeErr != nil {
			err = closeErr
			return
		}
	}()

	err = png.Encode(out, newImg)
	if err != nil {
		return
	}

	return
}

// ReplaceFace 替换人脸
func ReplaceFace(sourceImgUrl string, targetImgUrl string, outputImgUrl string) (err error) {
	outputImgUrl, _ = filepath.Abs(outputImgUrl)
	err = os.MkdirAll(filepath.Dir(outputImgUrl), os.ModePerm)
	if err != nil {
		fmt.Println(err, outputImgUrl)
		return
	}

	var msg []byte
	switch runtime.GOOS {
	case "windows":
		command := `python D:\Project\demo\roop\run.py -o %s -s %s -t %s`
		command = fmt.Sprintf(command, outputImgUrl, sourceImgUrl, targetImgUrl)
		msg, err = exec.Command("powershell", command).Output()
	default:
		log.Fatalln("I don't support other os")
	}
	if err != nil {
		fmt.Println(err, outputImgUrl)
		return
	}

	if string(msg) != "" {
		fmt.Println(string(msg), outputImgUrl)
	}

	return
}
