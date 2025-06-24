package webserve

import (
	"github.com/gin-gonic/gin"
	"goTT/webserve/hello"
	"goTT/webserve/user"
)

func Serve() chan error {
	r := gin.Default()

	// Koppla endpoints från olika paket
	r.GET("/hello", hello.Handler)
	r.GET("/user/:name", user.Handler)

	errChan := make(chan error, 1)

	go func() {
		errChan <- r.Run("localhost:8080") // will send nil or error
	}()

	return errChan
}
