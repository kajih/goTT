package mcu

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Dummy in-memory data store
var computers = []Computer{
	{Id: 1, Type: "pc", Name: "Office PC", Ip: "192.168.1.10", Uptime: "3d 2h", Online: true},
	{Id: 2, Type: "mcu", Name: "Design Mac", Ip: "192.168.1.11", Uptime: "12d 8h", Online: true},
	{Id: 3, Type: "skull", Name: "Pi Node", Ip: "192.168.1.42", Uptime: "99d 22h", Online: true},
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/computers", handler)
}

func handler(c *gin.Context) {
	c.JSON(http.StatusOK, computers)
}
