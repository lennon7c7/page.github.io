package proxy_test

import (
	"encoding/json"
	"fmt"
	"page.github.io/pkg/log"
	"page.github.io/pkg/proxy"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "start", "----------")

	m.Run()

	fmt.Println("----------", time.Now().Format("2006-01-02 15:04:05"), "end", "----------")
}

// go test -timeout 0 -v pkg/proxy/proxy_test.go -run TestGet
func TestGet(t *testing.T) {
	webUrl := "https://api.my-ip.io/ip"
	//webUrl = "https://vip.titan007.com/changeDetail/handicap.aspx?id=2183603&companyID=3&l=0"
	proxyUrl := "socks5://127.0.0.1:1080"

	proxyUrl = "socks5://127.0.0.1:9980"
	data := proxy.Get(webUrl, proxyUrl)
	fmt.Println(data)

	proxyUrl = "socks5://127.0.0.1:9981"
	data = proxy.Get(webUrl, proxyUrl)
	fmt.Println(data)

	proxyUrl = "socks5://127.0.0.1:9982"
	data = proxy.Get(webUrl, proxyUrl)
	fmt.Println(data)

	proxyUrl = "socks5://127.0.0.1:9983"
	data = proxy.Get(webUrl, proxyUrl)
	fmt.Println(data)
}

// go test -timeout 0 -v pkg/proxy/proxy_test.go -run TestOpenaiCompletions
func TestOpenaiCompletions(t *testing.T) {
	type OpenaiCompletionsResult struct {
		Error struct {
			Message string      `json:"message"`
			Type    string      `json:"type"`
			Param   interface{} `json:"param"`
			Code    interface{} `json:"code"`
		} `json:"error"`

		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Text         string      `json:"text"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	apiKey := "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	apiUrl := "https://api.openai.com/v1/completions"
	postHead := [][]string{
		{"Content-Type", "application/json"},
		{"Authorization", "Bearer " + apiKey},
	}
	postData := `{
	"model": "text-davinci-003",
  	"max_tokens": 300,
    "prompt": "List slippers attributes can see",
    "temperature": 0.9,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0.6,
    "stop": [
        " Human:",
        " AI:"
    ]
}`
	body, err := proxy.Post(apiUrl, postData, postHead)
	if err != nil {
		log.Error(err)
		return
	}

	var result OpenaiCompletionsResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err)
		return
	}

	if result.Error.Message != "" {
		log.Error(result.Error.Message)
		return
	}

	fmt.Println(result.Choices[0].Text)
}
