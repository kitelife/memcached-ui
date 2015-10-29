package main

import (
	"flag"
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

	if appConfig.Basic_auth.On == "yes" {
		r.Use(gin.BasicAuth(gin.Accounts{appConfig.Basic_auth.Username: appConfig.Basic_auth.Password}))
	}

	r.GET("/", controller.Home)
	r.POST("/do", controller.Do)

	var port string
	flag.StringVar(&port, "listen", ":8080", "listen address")
	flag.Parse()
	r.Run(port)
}
