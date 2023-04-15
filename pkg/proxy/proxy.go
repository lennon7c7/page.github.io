package proxy

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"page.github.io/pkg/log"
	"strings"
)

var DefaultUrl = "socks5://127.0.0.1:1080"

func GetUrl() (uri *url.URL) {
	uri, err := url.Parse(DefaultUrl)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func Get(webUrl string, proxyUrl string) (data string) {
	uri, err := url.Parse(proxyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}

	resp, err := client.Get(webUrl)
	if err != nil {
		log.Error(err)
		return
	}
	defer func(Body io.ReadCloser) {
		curlErr := Body.Close()
		if curlErr != nil {
			err = curlErr
			log.Error(err)
		}
	}(resp.Body)
	byteData, _ := io.ReadAll(resp.Body)

	data = string(byteData)

	return
}

func GetHtmlDom(webUrl string) (document *goquery.Document, err error) {
	uri, err := url.Parse(DefaultUrl)
	if err != nil {
		return
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}

	resp, err := client.Get(webUrl)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		curlErr := Body.Close()
		if curlErr != nil {
			err = curlErr
		}
	}(resp.Body)
	byteData, _ := io.ReadAll(resp.Body)

	document, err = goquery.NewDocumentFromReader(strings.NewReader(string(byteData)))
	if err != nil {
		return
	}

	return
}

func Post(apiUrl string, postData string, postHead [][]string) (body []byte, err error) {
	uri, err := url.Parse(DefaultUrl)
	if err != nil {
		log.Error(err)
		return
	}

	method := "POST"

	payload := strings.NewReader(postData)

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}

	req, err := http.NewRequest(method, apiUrl, payload)
	if err != nil {
		return
	}

	for _, i2 := range postHead {
		req.Header.Add(i2[0], i2[1])
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		curlErr := Body.Close()
		if curlErr != nil {
			err = curlErr
			log.Error(err)
		}
	}(res.Body)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}
