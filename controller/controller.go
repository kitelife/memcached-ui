package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/memcached"
)

type StatsInfoStruct struct {
	Server string
	Pid string
	Version string
	Uptime string
	MaxMemoryLimit string
	CurrMemoryUsage string
	CurrItems string
	CurrConnections string
	GetHits string
	GetMisses string
}

func getAppConfig(c *gin.Context) config.AppConfigStruct {
	appConf, _ := c.Get("app_conf")
	return appConf.(config.AppConfigStruct)
}

func getStatsInfo(server string)(map[string]string, error) {
	serverParts := strings.Split(server, ":")
	host := serverParts[0]
	port, _ := strconv.Atoi(serverParts[1])

	m := memcached.Memcached{}
	err := m.New(host, port)
	if err != nil {
		return nil, err
	}
	defer m.Close()
 	return m.Stats()
}

func checkValidServer(targetServer string, serverList []string)bool {
	for _, server := range serverList {
		if strings.Compare(server, targetServer) == 0 {
			return true
		}
	}
	return false
}

func formatUptime(uptime int) string {
	day := uptime / 86400
	hour := uptime % 86400 / 3600
	minute := uptime % 3600 / 60
	second := uptime % 60
	return fmt.Sprintf("%d天%d时%d分%d秒", day, hour, minute, second)
}

func formatMemoryUsage(usageBytes int) string {
	usageMB := usageBytes / 1024 / 1024
	return fmt.Sprintf("%dMB", usageMB)
}

func statsMap2Struct(statsMapper map[string]string) StatsInfoStruct {
	uptime, _ := strconv.Atoi(statsMapper["uptime"])
	maxMemoryLimit, _ := strconv.Atoi(statsMapper["limit_maxbytes"])
	currMemoryUsage, _ := strconv.Atoi(statsMapper["bytes"])

	return StatsInfoStruct {
		Pid: statsMapper["pid"],
		Version: statsMapper["version"],
		Uptime: formatUptime(uptime),
		MaxMemoryLimit: formatMemoryUsage(maxMemoryLimit),
		CurrMemoryUsage: formatMemoryUsage(currMemoryUsage),
		CurrItems: statsMapper["curr_items"],
		CurrConnections: statsMapper["curr_connections"],
		GetHits: statsMapper["get_hits"],
		GetMisses: statsMapper["get_misses"],
	}
}

func Home(c *gin.Context) {
	ac := getAppConfig(c)

	hostPortList := make([]string, 0, 100)
	for _, v := range ac.Servers {
		hostPortList = append(hostPortList, v)
	}

	targetServer := c.Query("server")
	if !checkValidServer(targetServer, hostPortList) {
		targetServer = hostPortList[0]
	}
	statsInfo, _ := getStatsInfo(targetServer)
	structedStatsInfo := statsMap2Struct(statsInfo)
	structedStatsInfo.Server = targetServer

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Servers": ac.Servers,
		"StatsInfo": structedStatsInfo,
	})
}
