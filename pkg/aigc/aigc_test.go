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
	definedCheckpoint := "chilloutmix.safetensors [fc2511737a]"

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
		options, err := aigc.GetOptions()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if options.SdModelCheckpoint != definedCheckpoint {
			fmt.Printf("options.SdModelCheckpoint != definedCheckpoint: %v\n", options.SdModelCheckpoint != definedCheckpoint)
			err = aigc.PostOptions(aigc.OptionsResponse{
				SdModelCheckpoint: definedCheckpoint,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		outputFilename := pathName + "/" + fmt.Sprintf("%08d", seed) + ".jpg"
		aigc.Txt2img(prompt, outputFilename, steps, seed)

		// restore
		if options.SdModelCheckpoint != definedCheckpoint {
			fmt.Printf("options.SdModelCheckpoint != definedCheckpoint: %v\n", options.SdModelCheckpoint != definedCheckpoint)
			err = aigc.PostOptions(options)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
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

func TestApiSegmentAnything(t *testing.T) {
	inputImgFile := "https://segment-anything.com/assets/gallery/GettyImages-1191014275.jpg"
	outputImgBase64, err := img.Http2Base64(inputImgFile)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(aigc.ApiSegmentAnything(outputImgBase64))
}

func TestGenerateMask(t *testing.T) {
	inputImgFile := "https://segment-anything.com/assets/gallery/GettyImages-1191014275.jpg"
	outputImgBase64, err := img.Http2Base64(inputImgFile)
	if err != nil {
		return
	}

	t.Log(aigc.GenerateMask(outputImgBase64))
}
