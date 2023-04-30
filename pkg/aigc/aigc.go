package aigc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/log"
	"strings"
	"time"
)

var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"

type Txt2ImgRequest struct {
	EnableHr          bool     `json:"enable_hr"`
	DenoisingStrength int      `json:"denoising_strength"`
	FirstphaseWidth   int      `json:"firstphase_width"`
	FirstphaseHeight  int      `json:"firstphase_height"`
	HrScale           int      `json:"hr_scale"`
	HrUpscaler        string   `json:"hr_upscaler"`
	HrSecondPassSteps int      `json:"hr_second_pass_steps"`
	HrResizeX         int      `json:"hr_resize_x"`
	HrResizeY         int      `json:"hr_resize_y"`
	Prompt            string   `json:"prompt"`
	Styles            []string `json:"styles"`
	Seed              int      `json:"seed"`
	Subseed           int      `json:"subseed"`
	SubseedStrength   int      `json:"subseed_strength"`
	SeedResizeFromH   int      `json:"seed_resize_from_h"`
	SeedResizeFromW   int      `json:"seed_resize_from_w"`
	SamplerName       string   `json:"sampler_name"`
	BatchSize         int      `json:"batch_size"`
	NIter             int      `json:"n_iter"`
	Steps             int      `json:"steps"`
	CfgScale          int      `json:"cfg_scale"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	RestoreFaces      bool     `json:"restore_faces"`
	Tiling            bool     `json:"tiling"`
	DoNotSaveSamples  bool     `json:"do_not_save_samples"`
	DoNotSaveGrid     bool     `json:"do_not_save_grid"`
	NegativePrompt    string   `json:"negative_prompt"`
	Eta               int      `json:"eta"`
	SChurn            int      `json:"s_churn"`
	STmax             int      `json:"s_tmax"`
	STmin             int      `json:"s_tmin"`
	SNoise            int      `json:"s_noise"`
	OverrideSettings  struct {
	} `json:"override_settings"`
	OverrideSettingsRestoreAfterwards bool          `json:"override_settings_restore_afterwards"`
	ScriptArgs                        []interface{} `json:"script_args"`
	SamplerIndex                      string        `json:"sampler_index"`
	ScriptName                        string        `json:"script_name"`
	SendImages                        bool          `json:"send_images"`
	SaveImages                        bool          `json:"save_images"`
	AlwaysonScripts                   struct {
	} `json:"alwayson_scripts"`
}

type Txt2ImgResponse struct {
	Images     []string `json:"images"`
	Parameters struct{} `json:"parameters"`
	Info       string   `json:"info"`
}

type PngContext struct {
	Parameters string
}

func Txt2img() {
	apiUrl := "http://127.0.0.1:7860/sdapi/v1/txt2img"
	method := "POST"

	payload := strings.NewReader(`{
	 "sd_model_checkpoint": "chilloutmix.safetensors [fc2511737a]",
	 "prompt": "<lora:koreanDollLikeness_v15:0.7>, masterpiece, best quality, ((((1girl)))), ((((huge breasts, detail breasts)))), ((((side-tie_bikini)))), ((((looking at viewer)))), ((((closeup)))), ((((detail arms, arms behind head)))), light blush",
	 "seed": -1,
	 "subseed": -1,
	 "batch_size": 1,
	 "steps": 20,
	 "cfg_scale": 7,
	 "width": 500,
  	 "height": 900,
	 "restore_faces": false,
	 "tiling": false,
	 "do_not_save_samples": false,
	 "do_not_save_grid": false,
	 "negative_prompt": "lowres, bad anatomy, ((((bad hands)))), bad feet, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry",
	 "eta": 0,
	 "s_churn": 0,
	 "s_tmax": 0,
	 "s_tmin": 0,
	 "s_noise": 1,
	 "override_settings": {},
	 "override_settings_restore_afterwards": true,
	 "script_args": [],
	 "sampler_index": "Euler",
	 "script_name": "",
	 "send_images": true,
	 "save_images": false,
	 "alwayson_scripts": {}
	}`)

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

	err = os.MkdirAll(BaseDownloadImgPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i2 := range resultData.Images {
		filePath := BaseDownloadImgPath + time.Now().Format("20060102150405") + "example.jpg"
		ddd, err := base64.StdEncoding.DecodeString(i2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = os.WriteFile(filePath, ddd, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(filePath)
	}
}
