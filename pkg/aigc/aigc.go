package aigc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"page.github.io/pkg/file"
	"page.github.io/pkg/log"
	"path"
	"path/filepath"
	"strings"
)

var BaseDownloadImgPath = "../../images/" + file.GetNameWithoutExt() + "/"

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
}

type PngContext struct {
	Parameters string
}

func Txt2img(outputFilename string) {
	apiUrl := "http://127.0.0.1:7860/sdapi/v1/txt2img"
	method := "POST"

	txt2ImgRequest := Txt2ImgRequest{
		SdModelCheckpoint: "chilloutmix.safetensors [fc2511737a]",
		Prompt:            "<lora:koreanDollLikeness_v15:0.7>, masterpiece, best quality, ((((1girl)))), ((((huge breasts, detail breasts)))), ((((looking at viewer)))), ((((closeup)))), ((((detail arms, arms behind head)))), light blush, ((((sexy lingerie))))",
		Seed:              -1,
		Subseed:           -1,
		BatchSize:         1,
		Steps:             30,
		CfgScale:          7,
		Width:             500,
		Height:            900,
		RestoreFaces:      false,
		Tiling:            false,
		DoNotSaveSamples:  false,
		DoNotSaveGrid:     false,
		NegativePrompt:    "lowres, bad anatomy, ((((bad hands)))), bad feet, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry",
		Eta:               0,
		SChurn:            0,
		STmax:             0,
		STmin:             0,
		SNoise:            1,
		//OverrideSettings:                  struct{}{},
		OverrideSettingsRestoreAfterwards: true,
		//ScriptArgs:                        nil,
		SamplerIndex: "Euler a",
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
		fmt.Println(outputFilename)
	}
}
