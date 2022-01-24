package ctx

import (
	"fmt"
	"github.com/hyperq/jpkg/db/mongo"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
)

var defaultFormat = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	statusColor = param.StatusCodeColor()
	methodColor = param.MethodColor()
	resetColor = param.ResetColor()

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf(
		"[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

var conf = gin.LoggerConfig{Formatter: defaultFormat}

const ErrorTypePrivate gin.ErrorType = 1 << 0

// Logger 日志
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(conf)
}

var mongolog = mongo.NewLog("access")

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter

	out := os.Stdout

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.Request.Header.Get("X-Real-IP")
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		logb, err := jsoniter.Marshal(
			map[string]interface{}{
				"ip":      param.ClientIP,
				"method":  param.Method,
				"code":    param.StatusCode,
				"emsg":    param.ErrorMessage,
				"path":    param.Path,
				"latency": param.Latency.Milliseconds(),
				"atime":   param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		_, _ = mongolog.Write(logb)
		_, _ = fmt.Fprint(out, formatter(param))
	}
}
