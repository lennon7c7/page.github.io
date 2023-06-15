package middleware

import (
	"github.com/gin-gonic/gin"
)

func KV(key string, value interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
		c.Next()
	}
}
