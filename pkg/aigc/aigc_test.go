package aigc_test

import (
	"fmt"
	"page.github.io/pkg/aigc"
	"page.github.io/pkg/img"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

func TestTxt2img(t *testing.T) {
	prompt := "wedding dress"
	prompt = "lucency full dress"
	prompt = "cropped jacket"
	prompt = "revealing dress"
	prompt = "v-neck shirts"
	//prompt = "Sailor dress"
	prompt = "hoodie"
	//prompt = "robe"
	//prompt = "cape"
	steps := 30
	pathName := aigc.BaseDownloadImgPath + strings.ReplaceAll(prompt+" steps "+fmt.Sprintf("%d", steps), " ", "-")

	minSeed := img.GetMaxFilename(pathName)
	if minSeed > 0 {
		minSeed++
	}
	fmt.Printf("minSeed: %v\n", minSeed)
	for seed := minSeed; seed < 99999999; seed++ {
		outputFilename := pathName + "/" + fmt.Sprintf("%08d", seed) + ".jpg"
		aigc.Txt2img(prompt, outputFilename, steps, seed)
	}
}

func TestImgRemoveBackgroundByUrl(t *testing.T) {
	inputImgUrl := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"
	outputImgBase64, err := aigc.ImgRemoveBackgroundByUrl(inputImgUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(outputImgBase64)
}
