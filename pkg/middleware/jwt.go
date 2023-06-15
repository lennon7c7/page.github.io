package middleware

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"page.github.io/pkg/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrorToken = "非法token，无权访问"
var ErrorTokenExpired = "令牌过期，请重新获取"

// JWTCheck 通过请求头验证签名
func JWTCheck(publicKey *rsa.PublicKey, skipPaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.RequestURI
		token := c.Request.Header.Get(util.JWTXToken)
		for _, skipPath := range skipPaths {
			// 特殊情况：当某些接口属于非强制携带TOKEN，又在过滤名单时，允许往下执行
			if token == "" && strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		if token == "" {
			util.Resp(c, util.CodeErrorToken, ErrorToken, "")
			return
		}

		tokenClaims, err := util.JWTCheck(publicKey, token)

		if err != nil {
			util.Resp(c, util.CodeErrorToken, ErrorToken, "")
			return
		}

		if !tokenClaims.Valid {
			util.Resp(c, util.CodeErrorTokenExpired, ErrorTokenExpired, "")
			return
		}

		claims := tokenClaims.Claims
		userID, _ := strconv.ParseInt(claims.(*jwt.StandardClaims).Issuer, 10, 0)

		c.Set(util.JWTClaims, claims)
		c.Set(util.JWTUserID, userID)
		c.Next()
	}
}
