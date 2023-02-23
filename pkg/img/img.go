package img

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"page.github.io/pkg/array"
	"page.github.io/pkg/file"
	"path"
)

var ExtList = []string{".jpg", ".jpeg", ".png", ".webp"}

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

func GetFiles(pathName string) (files []string) {
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

func Cut(inputFile string, outputFile string, outputWidth int, outputHeight int) {
	rectangle := image.Rect(0, 0, outputWidth, outputHeight)
	rgba := image.NewRGBA(rectangle)

	img, _ := os.Open(inputFile)
	defer func(img *os.File) {
		err := img.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(img)
	img1, _ := jpeg.Decode(img)

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
		fmt.Println(err)
		return
	}

	err = jpeg.Encode(tempFile, rgba, &jpeg.Options{Quality: 70})
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
