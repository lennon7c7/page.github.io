package aigc_test

import (
	"fmt"
	"page.github.io/pkg/aigc"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -timeout 0 -v pkg/aigc/aigc_test.go -run TestTxt2img
func TestTxt2img(t *testing.T) {
	for seed := 0; seed < 99999999; seed++ {
		steps := 10
		tag := "sexy lingerie steps " + fmt.Sprintf("%d", steps)
		outputFilename := aigc.BaseDownloadImgPath + strings.ReplaceAll(tag, " ", "-") + "/" + fmt.Sprintf("%08d", seed) + ".jpg"
		aigc.Txt2img(outputFilename, steps, seed)
	}
}
