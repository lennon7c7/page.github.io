package aigc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/img"
	"page.github.io/pkg/log"
	"page.github.io/pkg/util"
	"path"
	"path/filepath"
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

func GetOptions() (options OptionsResponse, err error) {
	apiUrl := UrlStableDiffusion + "sdapi/v1/options"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, nil)

	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &options)
	if err != nil {
		return
	}

	return
}

func PostOptions(options OptionsResponse) (err error) {
	apiUrl := UrlStableDiffusion + "sdapi/v1/options"
	method := "POST"

	optionsRequest := options
	newBuffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(newBuffer)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(optionsRequest)
	if err != nil {
		return
	}
	payload := strings.NewReader(newBuffer.String())

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, payload)
	if err != nil {
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(res.Body)

	return
}

func Txt2img(prompt string, outputFilename string, steps int, seed int) {
	if file.Exists(outputFilename) {
		//log.Info("文件已存在，跳过")
		return
	}

	apiUrl := UrlStableDiffusion + "sdapi/v1/txt2img"
	method := "POST"

	txt2ImgRequest := Txt2ImgRequest{
		Prompt: "((((sfw)))), <lora:koreanDollLikeness_v15:0.7>, masterpiece, best quality, ((((1girl)))), ((((huge breasts, detail breasts)))), ((((looking at viewer)))), ((((closeup)))), ((((detail arms)))), light blush, ((((" + prompt + "))))",
		//Seed:              -1,
		Seed: seed,
		//Subseed:           -1,
		BatchSize:        1,
		Steps:            steps,
		CfgScale:         7,
		Width:            480,
		Height:           880,
		RestoreFaces:     false,
		Tiling:           false,
		DoNotSaveSamples: false,
		DoNotSaveGrid:    false,
		NegativePrompt:   "((((nsfw)))), lowres, bad anatomy, ((((bad hands)))), bad feet, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry",
		Eta:              0,
		SChurn:           0,
		STmax:            0,
		STmin:            0,
		SNoise:           1,
		//OverrideSettings:                  struct{}{},
		OverrideSettingsRestoreAfterwards: true,
		//ScriptArgs:                        nil,
		//SamplerIndex: "Euler a",
		// DPM adpative最快成型，但后面就缺乏变化，感觉适合快速试验提示词组合
		SamplerIndex: "DPM++ SDE Karras",
		//ScriptName:                        "",
		SendImages: true,
		SaveImages: false,
		//AlwaysonScripts:                   struct{}{},
	}

	newBuffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(newBuffer)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(txt2ImgRequest)
	if err != nil {
		fmt.Println("Unable to convert the struct to a JSON string")
		return
	}
	payload := strings.NewReader(newBuffer.String())

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = closeErr
			log.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var resultData Txt2ImgResponse
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		fmt.Println(err)
		return
	}

	filePath := path.Dir(outputFilename)
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i2 := range resultData.Images {
		ddd, err := base64.StdEncoding.DecodeString(i2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = os.WriteFile(outputFilename, ddd, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}

		outputFilename, _ = filepath.Abs(outputFilename)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), outputFilename)
	}
}

func Img2img(request Img2ImgRequest) (response Txt2ImgResponse, err error) {
	apiUrl := UrlStableDiffusion + "sdapi/v1/img2img"
	method := "POST"

	newBuffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(newBuffer)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(request)
	if err != nil {
		return
	}
	payload := strings.NewReader(newBuffer.String())

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, payload)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Detail != "" {
		err = errors.New(response.Detail)
		return
	}

	for i, image := range response.Images {
		response.Images[i] = "data:image/jpeg;base64," + image
	}

	return
}

func ImgRemoveBackgroundByBase64(inputImgBase64 string) (outputImgBase64 string, err error) {
	apiUrl := "http://192.168.31.238:7860/sdapi/v1/extra-single-image"
	method := "POST"

	payload := strings.NewReader(`{
  "resize_mode": 0,
  "show_extras_results": true,
  "gfpgan_visibility": 0,
  "codeformer_visibility": 0,
  "codeformer_weight": 0,
  "upscaling_resize": 2,
  "upscaling_resize_w": 512,
  "upscaling_resize_h": 512,
  "upscaling_crop": true,
  "upscaler_1": "None",
  "upscaler_2": "None",
  "extras_upscaler_2_visibility": 0,
  "upscale_first": false,
  "image": "` + inputImgBase64 + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, payload)
	if err != nil {
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	type Response struct {
		// 正常返回
		HtmlInfo string `json:"html_info"`
		Image    string `json:"image"`

		// 错误返回
		Error  string `json:"error"`
		Detail string `json:"detail"`
		Body   string `json:"body"`
		Errors string `json:"errors"`
	}

	var resultData Response
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		return
	}

	if resultData.Detail != "" {
		err = errors.New(resultData.Detail)
		return
	}

	outputImgBase64 = resultData.Image
	return
}

func ImgRemoveBackgroundByUrl(inputImgUrl string) (outputImgBase64 string, err error) {
	inputImgBase64, err := img.Url2Base64(inputImgUrl)

	outputImgBase64, err = ImgRemoveBackgroundByBase64(inputImgBase64)

	return
}

func RouterImg2Img(c *gin.Context) {
	var request Img2ImgRequest
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	response, err := Img2img(request)
	if err != nil {
		util.ErrorBusiness(c, err.Error())
		return
	}

	util.OKData(c, gin.H{
		"images": response.Images,
	})
}
