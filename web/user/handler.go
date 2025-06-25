package user

import (
	"github.com/gin-gonic/gin"
)

func Handler(context *gin.Context) {
	context.String(200, "Hello, World!")
}
