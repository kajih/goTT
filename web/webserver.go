package web

import (
	"github.com/gin-gonic/gin"
	"goTT/web/hello"
	"goTT/web/user"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.GET("/v0/hello", hello.Handler)
	r.GET("/v0/user/:name", user.Handler)
	r.Static("/web", "./web_static")

	// (Optional) fallback for SPA routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./web_static/index.html")
	})

	return r
}
