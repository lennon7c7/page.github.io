package main

import (
	"flag"
	"fmt"
	"page.github.io/pkg/aigc"
	"page.github.io/pkg/img"
	"strings"
)

var prompt = flag.String("prompt", "v-neck shirts", "prompt")
var steps = flag.Int("steps", 30, "steps")

// go run .\cmd\sd\main.go
// go build -o pkg\aigc\sd.exe .\cmd\sd\main.go
// ./sd.exe -prompt="school uniform" -steps=1
func main() {
	flag.Parse()

	pathName := aigc.BaseDownloadImgPath + strings.ReplaceAll(*prompt+" steps "+fmt.Sprintf("%d", *steps), " ", "-")

	minSeed := img.GetMaxFilename(pathName)
	if minSeed > 0 {
		minSeed++
	}
	fmt.Printf("minSeed: %v\n", minSeed)
	for seed := minSeed; seed < 99999999; seed++ {
		outputFilename := pathName + "/" + fmt.Sprintf("%08d", seed) + ".jpg"
		aigc.Txt2img(*prompt, outputFilename, *steps, seed)
	}
}
