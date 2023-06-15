package router

import (
	"github.com/gin-gonic/gin"
	"page.github.io/pkg/aigc"
	"page.github.io/pkg/config"
)

func NewRouter(router *gin.Engine, config *config.Config) {
	router.StaticFile("favicon.ico", "./favicon.ico")
	router.StaticFile("/", "./index.html")
	for k, v := range config.Static {
		router.StaticFile(k, v)
	}
	router.Static("/static", "./static")
	router.Static(config.File.WebRelativePath, config.File.WebUploadRoot) //文件访问

	router.POST("img2img", aigc.RouterImg2Img)
	router.POST("generate-mask", aigc.RouterGenerateMask)
	router.POST("remove-img-background", aigc.RouterImgRemoveBackgroundByBase64)
}
