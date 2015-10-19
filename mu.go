package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/controller"
)

const (
	MU_CONFIG_PATH = "./app.json"
)

func muConfigMiddleware(conf config.MUConfigStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mu_conf", conf)
		c.Next()
	}
}

func main() {
	muConfig, err := config.LoadMUConfig(MU_CONFIG_PATH)
	if err != nil {
		fmt.Println("发生错误：", err.Error())
		os.Exit(-1)
	}

	r := gin.Default()
	r.Static("/assets", "./ui/assets")
	r.LoadHTMLGlob("ui/templates/*")
	r.Use(muConfigMiddleware(muConfig))

	r.GET("/", controller.Home)
	r.POST("/do", controller.Do)
	r.Run(":8080")
}
