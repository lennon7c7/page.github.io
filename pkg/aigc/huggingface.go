package aigc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"strings"
)

var HuggingFaceToken = ""

type HuggingFaceHuggingFaceObjectDetectionResponse struct {
	Score float64 `json:"score"`
	Label string  `json:"label"`
	Box   struct {
		Xmin int `json:"xmin"`
		Ymin int `json:"ymin"`
		Xmax int `json:"xmax"`
		Ymax int `json:"ymax"`
	} `json:"box"`
}

// HuggingFaceImg2TagsByRecognizeAnythingModel Recognize Anything Model
// @demo https://huggingface.co/spaces/xinyu1205/Recognize_Anything-Tag2Text
func HuggingFaceImg2TagsByRecognizeAnythingModel(base64Img string) (englishTag string, chineseTag string, err error) {
	type Response struct {
		Msg    string `json:"msg"`
		Output struct {
			Data            []string `json:"data"`
			IsGenerating    bool     `json:"is_generating"`
			Duration        float64  `json:"duration"`
			AverageDuration float64  `json:"average_duration"`
		} `json:"output"`
		Success bool `json:"success"`
	}

	sessionHash := generateRandomString11Length()

	url := "wss://xinyu1205-recognize-anything-tag2text.hf.space/queue/join"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			err = closeErr
		}
	}()

	// 接收消息
	for {
		messageType, messageByte, err := conn.ReadMessage()
		if messageType == -1 {
			fmt.Printf("messageType == -1: %v\n", messageType == -1)
			return englishTag, chineseTag, err
		}

		if err != nil {
			return englishTag, chineseTag, err
		}

		var response Response
		err = json.Unmarshal(messageByte, &response)
		if err != nil {
			return englishTag, chineseTag, err
		}

		// 判断消息
		switch true {
		case "send_hash" == response.Msg:
			// 发送数据
			data := map[string]interface{}{
				"fn_index":     2,
				"session_hash": sessionHash,
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				return englishTag, chineseTag, err
			}
			err = conn.WriteMessage(websocket.TextMessage, dataBytes)
			if err != nil {
				return englishTag, chineseTag, err
			}

			break

		case "send_data" == response.Msg:
			// 发送数据
			data := map[string]interface{}{
				"data":         []string{base64Img},
				"event_data":   nil,
				"fn_index":     2,
				"session_hash": sessionHash,
			}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				return englishTag, chineseTag, err
			}
			err = conn.WriteMessage(websocket.TextMessage, dataBytes)
			if err != nil {
				return englishTag, chineseTag, err
			}

			break

		case "process_completed" == response.Msg:
			englishTag = response.Output.Data[0]

			if len(response.Output.Data) == 2 {
				chineseTag = response.Output.Data[1]
			}

			return englishTag, chineseTag, err
		}
	}
}

// HuggingFaceObjectDetection Object Detection
// @demo https://huggingface.co/facebook/detr-resnet-50
func HuggingFaceObjectDetection(base64Img string) (responses []HuggingFaceHuggingFaceObjectDetectionResponse, err error) {
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
	apiUrl := "https://api-inference.huggingface.co/models/facebook/detr-resnet-50"
	req, err := http.NewRequest("POST", apiUrl, buffer)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "image/png")

	if HuggingFaceToken != "" {
		req.Header.Add("Authorization", "Bearer "+HuggingFaceToken)
	}

	// 使用代理发送请求
	//client := http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(proxy.GetUrl()),
	//	},
	//}

	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
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
