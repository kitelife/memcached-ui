package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/picasso250/memcached-ui/config"
	"github.com/picasso250/memcached-ui/memcached"
	MiddlemanManager "github.com/picasso250/memcached-ui/middleman/manager"
	_ "github.com/picasso250/memcached-ui/middleman/middleman"
)

type StatsInfoStruct struct {
	InstanceID      string
	Source          string
	Pid             string
	Version         string
	Uptime          string
	MaxMemoryLimit  string
	CurrMemoryUsage string
	CurrItems       string
	CurrConnections string
	GetHits         string
	GetMisses       string
	GetRate         string
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

func getAppConfig(c *gin.Context) config.AppConfigStruct {
	appConf, _ := c.Get("app_conf")
	return appConf.(config.AppConfigStruct)
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

	GetHits :=         statsMapper["get_hits"]
	GetMisses :=       statsMapper["get_misses"]
	GetRate := "0"
	if len(GetHits) > 0 && len(GetMisses) > 0 {
		h, err := strconv.Atoi(GetHits)
		if err != nil {
			log.Fatal(err)
		}
		m, err := strconv.Atoi(GetMisses)
		if err != nil {
			log.Fatal(err)
		}
		GetRate = strconv.FormatFloat(float64(h) / float64(m+h) * 100, 'f', 1, 64)
	}
	return StatsInfoStruct{
		Pid:             statsMapper["pid"],
		Version:         statsMapper["version"],
		Uptime:          formatUptime(uptime),
		MaxMemoryLimit:  formatMemoryUsage(maxMemoryLimit),
		CurrMemoryUsage: formatMemoryUsage(currMemoryUsage),
		CurrItems:       statsMapper["curr_items"],
		CurrConnections: statsMapper["curr_connections"],
		GetHits:         GetHits,
		GetMisses:       GetMisses,
		GetRate:         GetRate,
	}
}

func Home(c *gin.Context) {
	ac := getAppConfig(c)

	hostPortList := make([]string, 0, 100)
	for k, _ := range ac.Instances {
		hostPortList = append(hostPortList, k)
	}

	instanceID := c.Query("instance")
	if _, ok := ac.Instances[instanceID]; ok == false {
		instanceID = hostPortList[0]
	}
	targetSource := ac.Instances[instanceID].Source

	infoErr := ""
	hasInfoErr := false
	statsInfo, err := getStatsInfo(targetSource)
	if err != nil {
		infoErr = err.Error()
		hasInfoErr = true
	}
	structedStatsInfo := statsMap2Struct(statsInfo)
	structedStatsInfo.InstanceID = instanceID
	structedStatsInfo.Source = targetSource

	c.HTML(http.StatusOK, "index.html", gin.H{
		"HasInfoErr": hasInfoErr,
		"InfoErr":    infoErr,
		"Instances":  ac.Instances,
		"StatsInfo":  structedStatsInfo,
	})
}

func Do(c *gin.Context) {
	ac := getAppConfig(c)

	instanceID := c.PostForm("instance")
	if _, ok := ac.Instances[instanceID]; ok == false {
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
	m, err := newMemcached(ac.Instances[instanceID].Source)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"msg":    "目标Memcached服务连接失败：" + err.Error(),
		})
		return
	}
	defer m.Close()

	targetInstanceConfig := ac.Instances[instanceID]
	targetMiddleman := MiddlemanManager.Get(targetInstanceConfig.MiddlemanName, targetInstanceConfig.MiddlemanConfig)
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
