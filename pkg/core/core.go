package core

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"os"
	"page.github.io/pkg/config"
	"page.github.io/pkg/middleware"
	"page.github.io/pkg/router"
	"strconv"
)

// @title Swagger API
// @version 1.0
// @description This is a swagger API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// Core @BasePath /
func Core() {
	config.CFG = config.NewConfig()

	if config.CFG.DB.DataSourceName != "" {
		// 初始化数据库
		config.DB, _ = config.NewDB(config.CFG)
	}

	if config.CFG.DB.DataSourceName != "" {
		// Casbin权限设置
		middleware.CasbinSetup(config.DB, config.CFG.Casbin.Model)
	}

	r := gin.Default()
	gin.SetMode(config.CFG.GinMode)

	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(config.CFG.JWT.PublicKey))

	// 设置中间件
	r.Use(middleware.Cors(), middleware.JWTCheck(publicKey, config.CFG.JWT.SkipPaths...))

	// 配置路由
	router.NewRouter(r, config.CFG)

	// swagger 设置
	if config.CFG.ShowSwagger {
		swaggerJson := ""
		swaggerJsons := []string{"./api/swagger.json", "../../api/swagger.json"}
		for _, v := range swaggerJsons {
			if _, err := os.Stat(v); err == nil {
				swaggerJson = v
				break
			}
		}

		if swaggerJson == "" {
			panic("查询 swagger.json 不存在，请重新尝试")
		}

		r.StaticFile("/swagger.json", swaggerJson) //文件访问
	}

	// 设定端口号
	if config.CFG.Https.On {
		r.Use(TlsHandler(config.CFG))
		_ = r.RunTLS(":"+strconv.Itoa(config.CFG.Https.Port), config.CFG.Https.Host+".crt", config.CFG.Https.Host+".key")
	} else {
		_ = r.Run(":" + strconv.Itoa(config.CFG.Http.Port))
	}
}

func TlsHandler(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     config.Https.Host + ":" + strconv.Itoa(config.Https.Port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
