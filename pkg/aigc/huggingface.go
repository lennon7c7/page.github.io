package aigc

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"page.github.io/pkg/log"
)

// Img2TagsByRecognizeAnythingModel Recognize Anything Model
// @demo https://huggingface.co/spaces/xinyu1205/Recognize_Anything-Tag2Text
func Img2TagsByRecognizeAnythingModel(base64Img string) (englishTag string, chineseTag string, err error) {
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
		log.Fatal("dial error:", err)
	}
	defer func(conn *websocket.Conn) {
		closeErr := conn.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(conn)

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
