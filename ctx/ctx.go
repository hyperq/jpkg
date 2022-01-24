package ctx

import (
	"github.com/gin-gonic/gin"
)

// Context ext content struct
type Context struct {
	*gin.Context
}

// HandlerFunc HandlerFunc
type HandlerFunc func(*Context)

// Handler ext gin.content
func H(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(Context)
		context.Context = c
		h(context)
	}
}
