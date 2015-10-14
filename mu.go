package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/controller"
)

const (
	APP_CONFIG_PATH = "./app.json"
)

func appConfigMiddleware(conf config.AppConfigStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_conf", conf)
		c.Next()
	}
}

func main() {
	appConfig, err := config.LoadAppConfig(APP_CONFIG_PATH)
	if err != nil {
		fmt.Println("发生错误：", err.Error())
		os.Exit(-1)
	}

	r := gin.Default()
	r.Static("/assets", "./ui/assets")
	r.LoadHTMLGlob("ui/templates/*")
	r.Use(appConfigMiddleware(appConfig))

	r.GET("/", controller.Home)
	r.Run(":8080")
}
