package main

import (
	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/controller"
)

func main() {
	r := gin.Default()
	r.Static("/assets", "./ui/assets")

	r.GET("/stats", controller.StatInfo)
	r.Run(":8080")
}
