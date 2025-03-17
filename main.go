package main

import (
	"github.com/gin-gonic/gin"
	"gocaptcha/internal/routers"
	"gocaptcha/internal/utils/server"
	"gocaptcha/pkg/config"
)

var debug = config.Config.GetBool("server.debug")
var port = ":" + config.Config.GetString("server.port")

func main() {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	server.Run(routers.R, port)
}
