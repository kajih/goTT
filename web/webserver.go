package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"goTT/web/hello"
	"goTT/web/mcu"
	"goTT/web/sse"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.Use(cors.Default())

	hello.RegisterRoutes(r)
	mcu.RegisterRoutes(r)
	sse.RegisterRoutes(r)

	r.Static("/web", "./web_static")

	// (Optional) fallback for SPA routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./web_static/index.html")
	})

	return r
}
