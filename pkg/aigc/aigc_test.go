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
	prompt := "cardigan"
	steps := 30
	pathName := aigc.BaseDownloadImgPath + strings.ReplaceAll(prompt+" steps "+fmt.Sprintf("%d", steps), " ", "-")

	minSeed := img.GetMaxFilename(pathName)
	if minSeed > 0 {
		minSeed++
	}
	for seed := minSeed; seed < 99999999; seed++ {
		outputFilename := pathName + "/" + fmt.Sprintf("%08d", seed) + ".jpg"
		aigc.Txt2img(prompt, outputFilename, steps, seed)
	}
}
