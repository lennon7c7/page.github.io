package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Log(logFile string) gin.HandlerFunc {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Now().Sub(startTime)

		log := zerolog.New(f).With().Timestamp().Logger()
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		log.Log().
			Str("method", c.Request.Method).
			Stringer("url", c.Request.URL).
			Str("uri", c.Request.RequestURI).
			Str("user-agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Dur("duration", duration).
			Str("ip", c.ClientIP()).
			Str("user", "").
			Msg("")
	}
}
