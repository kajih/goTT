package hello

import (
	"github.com/gin-gonic/gin"
)

func handler(context *gin.Context) {
	context.String(200, "Hello, World!")
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/hello", handler)
}
