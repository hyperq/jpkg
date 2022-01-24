package ctx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Context ext content struct
type Context struct {
	*gin.Context
}

// HandlerFunc HandlerFunc
type HandlerFunc func(*Context)

// H Handler ext gin.content
func H(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		h(context)
	}
}

// Cors 跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		Origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", Origin)                           // header的类型
		c.Header("Access-Control-Allow-Headers", "Content-Type,X-Requested-With") // header的类型
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")    // 允许post访问
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 处理请求
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next() // 处理请求
	}
}
