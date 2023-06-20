package aigc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"page.github.io/pkg/img"
	"strings"
	"time"
)

var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"
var UrlStableDiffusion = "http://127.0.0.1:7860/"

type OptionsResponse struct {
	SdModelCheckpoint string `json:"sd_model_checkpoint"`
}

type Txt2ImgRequest struct {
	SdModelCheckpoint string `json:"sd_model_checkpoint"`
	Prompt            string `json:"prompt"`
	Seed              int    `json:"seed"`
	Subseed           int    `json:"subseed"`
	BatchSize         int    `json:"batch_size"`
	Steps             int    `json:"steps"`
	CfgScale          int    `json:"cfg_scale"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	RestoreFaces      bool   `json:"restore_faces"`
	Tiling            bool   `json:"tiling"`
	DoNotSaveSamples  bool   `json:"do_not_save_samples"`
	DoNotSaveGrid     bool   `json:"do_not_save_grid"`
	NegativePrompt    string `json:"negative_prompt"`
	Eta               int    `json:"eta"`
	SChurn            int    `json:"s_churn"`
	STmax             int    `json:"s_tmax"`
	STmin             int    `json:"s_tmin"`
	SNoise            int    `json:"s_noise"`
	//OverrideSettings                  struct{}      `json:"override_settings"`
	OverrideSettingsRestoreAfterwards bool `json:"override_settings_restore_afterwards"`
	//ScriptArgs                        []interface{} `json:"script_args"`
	SamplerIndex string `json:"sampler_index"`
	ScriptName   string `json:"script_name"`
	SendImages   bool   `json:"send_images"`
	SaveImages   bool   `json:"save_images"`
	//AlwaysonScripts                   struct{}      `json:"alwayson_scripts"`
}

type Txt2ImgResponse struct {
	Images     []string `json:"images"`
	Parameters struct{} `json:"parameters"`
	Info       string   `json:"info"`
	Detail     string   `json:"detail"`
}

type Img2ImgRequest struct {
	Prompt                string   `json:"prompt"`
	NegativePrompt        string   `json:"negative_prompt"`
	Seed                  int64    `json:"seed"`
	BatchSize             int      `json:"batch_size"`
	NIter                 int      `json:"n_iter"`
	Steps                 int      `json:"steps"`
	CfgScale              int      `json:"cfg_scale"`
	Width                 int      `json:"width"`
	Height                int      `json:"height"`
	SamplerIndex          string   `json:"sampler_index"`
	ResizeMode            int      `json:"resize_mode"`
	DenoisingStrength     float64  `json:"denoising_strength"`
	MaskBlur              int      `json:"mask_blur"`
	InpaintingFill        int      `json:"inpainting_fill"`
	InpaintFullRes        bool     `json:"inpaint_full_res"`
	InpaintFullResPadding int      `json:"inpaint_full_res_padding"`
	InpaintingMaskInvert  int      `json:"inpainting_mask_invert"`
	InitImages            []string `json:"init_images"`
	Mask                  string   `json:"mask"`
}

type PngContext struct {
	Parameters string
}

type SegmentAnythingResponse struct {
	EncodedMask  string    `json:"encodedMask"`
	Bbox         []float64 `json:"bbox"`
	Score        float64   `json:"score"`
	PointCoord   []float64 `json:"point_coord"`
	UncertainIou float64   `json:"uncertain_iou"`
	Area         int       `json:"area"`
}

func GenerateMask(inputImgBase64 string) (outputImgBase64 string, err error) {
	imgWidth, imgHeight, err := img.GetImageSizeFromBase64(inputImgBase64)
	if err != nil {
		return
	}

	imgBase64RemoveBackground, err := ImgRemoveBackgroundByBase64(inputImgBase64)
	if err != nil {
		return
	}

	responses, err := ApiSegmentAnything(imgBase64RemoveBackground)
	if err != nil {
		return
	}

	if len(responses) < 2 {
		err = errors.New(`len(responses) < 2`)
		return
	}

	if len(responses[1].Bbox) < 4 {
		err = errors.New(`len(responses[1].Bbox) < 4`)
		return
	}
	bBox := responses[1].Bbox

	maskX, maskY, maskWidth, maskHeight := int(bBox[0]), int(bBox[1]), int(bBox[2]), int(bBox[3])
	outputImgBase64, err = img.GenerateRectMask(imgWidth, imgHeight, maskX, maskY, maskWidth, maskHeight)
	if err != nil {
		return
	}

	outputImgBase64 = "data:image/png;base64," + outputImgBase64

	return
}

// ApiSegmentAnything 细分任何内容
// @demo https://segment-anything.com/demo
func ApiSegmentAnything(base64Img string) (responses []SegmentAnythingResponse, err error) {
	substr := ","
	if strings.Contains(base64Img, substr) {
		// 兼容
		base64Img = strings.Split(base64Img, substr)[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return
	}

	// 将字节数组转换为缓冲区
	buffer := bytes.NewBuffer(decoded)

	// 创建一个 POST 请求
	apiUrl := "https://model-zoo.metademolab.com/predictions/automatic_masks"
	req, err := http.NewRequest("POST", apiUrl, buffer)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "image/jpg")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &responses)
	if err != nil {
		return
	}

	return
}

func GenerateMaskBySegmentAnything(inputBase64Img string) (cutImgList []string, maskList []string, err error) {
	originImgInterface, err := img.Base64ToImgInterface(inputBase64Img)
	if err != nil {
		return
	}

	responses, err := ApiSegmentAnything(inputBase64Img)
	if err != nil {
		return
	}

	for _, response := range responses {
		maskX, maskY, maskWidth, maskHeight := int(response.Bbox[0]), int(response.Bbox[1]), int(response.Bbox[2]), int(response.Bbox[3])
		_, base64Img, err := img.CutByBbox(originImgInterface, maskX, maskY, maskWidth, maskHeight)
		if err != nil {
			return nil, nil, err
		}

		tempMask, err := ImgRemoveBackgroundByBase64(base64Img)
		if err != nil {
			return nil, nil, err
		}

		outputFilename := generateRandomString11Length() + ".png"
		err = img.DrawTransparentBackground(tempMask, originImgInterface.Bounds().Dx(), originImgInterface.Bounds().Dy(), maskX, maskY, outputFilename)
		if err != nil {
			return nil, nil, err
		}

		mask, err := img.File2Base64(outputFilename)
		if err != nil {
			return nil, nil, err
		}
		_ = os.Remove(outputFilename)

		cutImgList = append(cutImgList, mask)

		mask, err = GenerateMaskByRembg(mask)
		if err != nil {
			return nil, nil, err
		}

		maskList = append(maskList, mask)
		break
	}

	return
}

func generateRandomString11Length() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 11)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
