package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/memcached"
	"net/http"
	"strconv"
	"strings"
)

func StatInfo(c *gin.Context) {
	server := strings.Split(c.Query("server"), ":")
	host := server[0]
	port, _ := strconv.Atoi(server[1])

	statType := c.Query("type")

	m := memcached.Memcached{}
	err := m.New(host, port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"msg":    err,
		})
		return
	}

	defer m.Close()

	mapper, err := m.Stats(statType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"msg":    err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mapper,
	})
}
