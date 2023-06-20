package aigc

import (
	"github.com/gin-gonic/gin"
	"page.github.io/pkg/img"
	"page.github.io/pkg/log"
	"page.github.io/pkg/util"
)

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

func RouterImgRemoveBackgroundByBase64(c *gin.Context) {
	type Request struct {
		Base64Img string `json:"base64Img"`
	}

	var request Request
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	response, err := ImgRemoveBackgroundByBase64(request.Base64Img)
	if err != nil {
		util.ErrorBusiness(c, err.Error())
		return
	}

	util.OKData(c, gin.H{
		"image": response,
	})
}

func RouterGenerateMask(c *gin.Context) {
	type Request struct {
		Base64Img string `json:"base64Img"`
	}

	var request Request
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	response, err := GenerateMask(request.Base64Img)
	if err != nil {
		util.ErrorBusiness(c, err.Error())
		log.Error(err)
		return
	}

	util.OKData(c, gin.H{
		"image": response,
	})
}

func RouterGenerateMaskByRembg(c *gin.Context) {
	type Request struct {
		Base64Img string `json:"base64Img"`
	}

	var request Request
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	response, err := GenerateMaskByRembg(request.Base64Img)
	if err != nil {
		util.ErrorBusiness(c, err.Error())
		return
	}

	util.OKData(c, gin.H{
		"image": response,
	})
}

func RouterGenerateMaskBySam(c *gin.Context) {
	type Request struct {
		Base64Img string `json:"base64Img"`
	}

	var request Request
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	imgWidth, imgHeight, err := img.GetImageSizeFromBase64(request.Base64Img)
	if err != nil {
		return
	}

	imgBase64RemoveBackground, err := ImgRemoveBackgroundByBase64(request.Base64Img)
	if err != nil {
		return
	}

	responses, err := ApiSegmentAnything(imgBase64RemoveBackground)
	if err != nil {
		return
	}

	var images []string
	for i, response := range responses {
		if i == 0 {
			continue
		}

		maskX, maskY, maskWidth, maskHeight := int(response.Bbox[0]), int(response.Bbox[1]), int(response.Bbox[2]), int(response.Bbox[3])
		image, err := img.GenerateRectMask(imgWidth, imgHeight, maskX, maskY, maskWidth, maskHeight)
		if err != nil {
			continue
		}

		image = "data:image/png;base64," + image
		images = append(images, image)
	}

	util.OKData(c, gin.H{
		"responses": responses,
		"images":    images,
	})
}

func RouterHuggingFaceImg2TagsByRecognizeAnythingModel(c *gin.Context) {
	type Request struct {
		Base64Img string `json:"base64Img"`
	}

	var request Request
	if err := c.ShouldBind(&request); err != nil {
		util.Error(c, "获取参数 错误，请重新尝试", err.Error())
		return
	}

	englishTag, chineseTag, err := HuggingFaceImg2TagsByRecognizeAnythingModel(request.Base64Img)
	if err != nil {
		util.ErrorBusiness(c, err.Error())
		return
	}

	util.OKData(c, gin.H{
		"englishTag": englishTag,
		"chineseTag": chineseTag,
	})
}
