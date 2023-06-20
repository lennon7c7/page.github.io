package aigc_test

import (
	"fmt"
	"page.github.io/pkg/aigc"
	"page.github.io/pkg/img"
	"page.github.io/pkg/log"
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
		log.Error(err)
		return
	}

	err = img.Base64ToFile(outputImgBase64, "../../images/output.png")
	if err != nil {
		log.Error(err)
		return
	}
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

func TestGenerateMaskBySegmentAnything(t *testing.T) {
	inputImgFile := "https://segment-anything.com/assets/gallery/GettyImages-1191014275.jpg"
	outputImgBase64, err := img.Http2Base64(inputImgFile)
	if err != nil {
		log.Error(err)
		return
	}

	imgList, maskList, err := aigc.GenerateMaskBySegmentAnything(outputImgBase64)
	if err != nil {
		log.Error(err)
		return
	}

	for i, base64 := range imgList {
		err = img.Base64ToFile(base64, "../../images/remove-background-"+fmt.Sprintf("%d", i)+".png")
		if err != nil {
			log.Error(err)
			continue
		}
	}

	for i, base64 := range maskList {
		err = img.Base64ToFile(base64, "../../images/mask-"+fmt.Sprintf("%d", i)+".png")
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func TestGenerateMask(t *testing.T) {
	inputImgFile := "https://segment-anything.com/assets/gallery/GettyImages-1191014275.jpg"
	outputImgBase64, err := img.Http2Base64(inputImgFile)
	if err != nil {
		log.Error(err)
		return
	}

	outputImgBase64, err = aigc.GenerateMask(outputImgBase64)
	if err != nil {
		log.Error(err)
		return
	}

	err = img.Base64ToFile(outputImgBase64, "../../images/output.png")
	if err != nil {
		log.Error(err)
		return
	}
}

func TestHuggingFaceImg2TagsByRecognizeAnythingModel(t *testing.T) {
	inputImgFile := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"
	outputImgBase64, err := img.Http2Base64(inputImgFile)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(aigc.HuggingFaceImg2TagsByRecognizeAnythingModel(outputImgBase64))
}

func TestHuggingFaceObjectDetection(t *testing.T) {
	//inputImgFile := "https://s3.amazonaws.com/a.storyblok.com/f/191576/1024x1024/f187412139/sample-03.png"
	//outputImgBase64, err := img.Http2Base64(inputImgFile)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	inputImgFile := "../../images/nature-1.jpg"
	outputImgBase64, err := img.File2Base64(inputImgFile)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(aigc.HuggingFaceObjectDetection(outputImgBase64))
}
