package controller

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/youngsterxyf/memcached-ui/config"
	"github.com/youngsterxyf/memcached-ui/memcached"
)

type StatsInfoStruct struct {
	Id              string
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
		if strings.Compare(targetAction, action) == 0 {
			return true
		}
	}
	return false
}

func getAppConfig(c *gin.Context) config.AppConfigStruct {
	appConf, _ := c.Get("app_conf")
	return appConf.(config.AppConfigStruct)
}

func genYiiKey(key string, yiiConf map[string]string) string {
	fmt.Println(key)
	innerKey := fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(yiiConf["app_name"]))) + key
	if yiiConf["hash"] == "yes" {
		innerKey = fmt.Sprintf("%x", md5.Sum([]byte(innerKey)))
	}
	fmt.Println(innerKey)
	return innerKey
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
	usageMB := usageBytes / 1024 / 1024
	return fmt.Sprintf("%dMB", usageMB)
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
	ac := getAppConfig(c)

	hostPortList := make([]string, 0, 100)
	for k, _ := range ac.Servers {
		hostPortList = append(hostPortList, k)
	}

	targetServer := c.Query("server")
	if _, ok := ac.Servers[targetServer]; ok == false {
		targetServer = hostPortList[0]
	}

	infoErr := ""
	hasInfoErr := false
	statsInfo, err := getStatsInfo(targetServer)
	if err != nil {
		infoErr = err.Error()
		hasInfoErr = true
	}
	structedStatsInfo := statsMap2Struct(statsInfo)
	structedStatsInfo.Server = targetServer
	structedStatsInfo.Id = ac.Servers[targetServer]

	c.HTML(http.StatusOK, "index.html", gin.H{
		"HasInfoErr": hasInfoErr,
		"InfoErr":    infoErr,
		"Servers":    ac.Servers,
		"StatsInfo":  structedStatsInfo,
	})
}

func Do(c *gin.Context) {
	ac := getAppConfig(c)

	targetServer := c.PostForm("server")
	if _, ok := ac.Servers[targetServer]; ok == false {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "不存在目标Memcached服务",
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
	m, err := newMemcached(targetServer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "目标Memcached服务连接失败：" + err.Error(),
		})
		return
	}
	defer m.Close()

	// 是否支持Yii使用缓存的方式
	useYii := false
	if ac.Yii["status"] == "on" {
		useYii = true
	}

	switch {
	case targetAction == "get":
		key := c.PostForm("key")
		if useYii {
			key = genYiiKey(key, ac.Yii)
		}
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
			"data":   string(resp),
		})
		return
	case targetAction == "set":
		key := c.PostForm("key")
		if useYii {
			key = genYiiKey(key, ac.Yii)
		}
		value := c.PostForm("value")
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
	case strings.Compare(targetAction, "delete") == 0:
		key := c.PostForm("key")
		if useYii {
			key = genYiiKey(key, ac.Yii)
		}
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
	case strings.Compare(targetAction, "flush_all") == 0:
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
