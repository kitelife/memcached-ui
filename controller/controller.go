package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/memcached"
	MiddlemanManager "github.com/youngsterxyf/memcached-ui/middleman/manager"
	_ "github.com/youngsterxyf/memcached-ui/middleman/middleman"
)

type StatsInfoStruct struct {
	App             string
	Server          string
	Pid             string
	Version         string
	Uptime          string
	MaxMemoryLimit  string
	CurrMemoryUsage string
	CurrItems       string
	CurrConnections string
	GetHits         string
	GetMisses       string
}

var actionAllowed []string = []string{"get", "set", "delete", "flush_all"}

func validAction(targetAction string) bool {
	for _, action := range actionAllowed {
		if targetAction == action {
			return true
		}
	}
	return false
}

func getMUConfig(c *gin.Context) config.MUConfigStruct {
	muConf, _ := c.Get("mu_conf")
	return muConf.(config.MUConfigStruct)
}

func newMemcached(server string) (memcached.Memcached, error) {
	serverParts := strings.Split(server, ":")
	host := serverParts[0]
	port, _ := strconv.Atoi(serverParts[1])

	m := memcached.Memcached{}
	err := m.New(host, port)
	return m, err
}

func getStatsInfo(server string) (map[string]string, error) {
	m, err := newMemcached(server)
	if err != nil {
		return nil, err
	}
	defer m.Close()
	return m.Stats()
}

func formatUptime(uptime int) string {
	day := uptime / 86400
	hour := uptime % 86400 / 3600
	minute := uptime % 3600 / 60
	second := uptime % 60
	return fmt.Sprintf("%d天%d时%d分%d秒", day, hour, minute, second)
}

func formatMemoryUsage(usageBytes int) string {
	usageKB := float32(usageBytes) / float32(1024)
	usageMB := usageKB / float32(1024)
	return fmt.Sprintf("%.2fMB (%.2fKB)", usageMB, usageKB)
}

func statsMap2Struct(statsMapper map[string]string) StatsInfoStruct {
	uptime, _ := strconv.Atoi(statsMapper["uptime"])
	maxMemoryLimit, _ := strconv.Atoi(statsMapper["limit_maxbytes"])
	currMemoryUsage, _ := strconv.Atoi(statsMapper["bytes"])

	return StatsInfoStruct{
		Pid:             statsMapper["pid"],
		Version:         statsMapper["version"],
		Uptime:          formatUptime(uptime),
		MaxMemoryLimit:  formatMemoryUsage(maxMemoryLimit),
		CurrMemoryUsage: formatMemoryUsage(currMemoryUsage),
		CurrItems:       statsMapper["curr_items"],
		CurrConnections: statsMapper["curr_connections"],
		GetHits:         statsMapper["get_hits"],
		GetMisses:       statsMapper["get_misses"],
	}
}

func Home(c *gin.Context) {
	muc := getMUConfig(c)

	hostPortList := make([]string, 0, 100)
	for k, _ := range muc.Servers {
		hostPortList = append(hostPortList, k)
	}

	targetApp := c.Query("app")
	if _, ok := muc.Servers[targetApp]; ok == false {
		targetApp = hostPortList[0]
	}
	targetServer := muc.Servers[targetApp].Source

	infoErr := ""
	hasInfoErr := false
	statsInfo, err := getStatsInfo(targetServer)
	if err != nil {
		infoErr = err.Error()
		hasInfoErr = true
	}
	structedStatsInfo := statsMap2Struct(statsInfo)
	structedStatsInfo.App = targetApp
	structedStatsInfo.Server = targetServer

	c.HTML(http.StatusOK, "index.html", gin.H{
		"HasInfoErr": hasInfoErr,
		"InfoErr":    infoErr,
		"Servers":    muc.Servers,
		"StatsInfo":  structedStatsInfo,
	})
}

func Do(c *gin.Context) {
	muc := getMUConfig(c)

	targetApp := c.PostForm("app")
	if _, ok := muc.Servers[targetApp]; ok == false {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "不存在目标应用",
		})
		return
	}
	targetAction := c.PostForm("action")
	if validAction(targetAction) == false {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "不存在目标action",
		})
		return
	}
	m, err := newMemcached(muc.Servers[targetApp].Source)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "目标Memcached服务连接失败：" + err.Error(),
		})
		return
	}
	defer m.Close()

	targetAppConfig := muc.Servers[targetApp]
	targetMiddleman := MiddlemanManager.Get(targetAppConfig.MiddlemanName, targetAppConfig.MiddlemanConfig)
	if targetMiddleman == nil {
		targetMiddleman = MiddlemanManager.Get("default", nil)
	}

	switch {
	case targetAction == "get":
		key := targetMiddleman.GenInnerKey(c.PostForm("key"))
		resp, err := m.Get(key)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"msg":    "获取缓存数据失败：" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   targetMiddleman.UnserializeValue(resp),
		})
		return
	case targetAction == "set":
		key := targetMiddleman.GenInnerKey(c.PostForm("key"))
		value := targetMiddleman.SerializeValue(c.PostForm("value"))
		expTime := c.DefaultPostForm("exp_time", "0")
		expTimeInt, err := strconv.Atoi(expTime)
		if err != nil {
			expTimeInt = 0
		}
		resp, err := m.Set(memcached.StorageCmdArgStruct{"key": key, "value": value, "expire_time": expTimeInt})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"msg":    "添加缓存失败：" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   string(resp),
		})
	case targetAction == "delete":
		key := targetMiddleman.GenInnerKey(c.PostForm("key"))
		resp, err := m.Delete(key)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"msg":    "删除缓存失败：" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   string(resp),
		})
	case targetAction == "flush_all":
		resp, err := m.FlushAll()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"msg":    "清空缓存失败：" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   string(resp),
		})
	}
}
