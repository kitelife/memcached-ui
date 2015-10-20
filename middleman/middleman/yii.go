package middleman

import (
	"crypto/md5"
	"fmt"
	"bufio"
	"strings"
	"encoding/json"
	"hash/crc32"

	"github.com/picasso250/memcached-ui/middleman/manager"
	"github.com/picasso250/memcached-ui/phpunserialize"
)

type YiiMiddleman struct {
	config map[string]string
}

func (ymm *YiiMiddleman) Config(config map[string]string) bool {
	ymm.config = config
	return true
}

func (ymm YiiMiddleman) GenInnerKey(key string) string {
	innerKey := fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(ymm.config["appName"]))) + key
	if ymm.config["hash"] == "yes" {
		innerKey = fmt.Sprintf("%x", md5.Sum([]byte(innerKey)))
	}
	return innerKey
}

func (ymm YiiMiddleman) SerializeValue(value string) string {
	return value
}

func (ymm YiiMiddleman) UnserializeValue(value string) interface{} {
	var data interface{}
	// 仅自动解析 PHP `serialize()`ed Array
	if len(value) > 2 && value[0] == 'a' && value[1] == ':' {
		data = phpunserialize.Parse(bufio.NewReader(strings.NewReader(value)))
		// Yii 模式下自动提取 JSON
		arr, ok := data.([]interface{})
		if (ok && len(arr) == 2) {
			mainstr, ok := arr[0].(string)
			if (ok && len(mainstr) >= 2 && (mainstr[0] == '{' || mainstr[0] == '[')) {
				err := json.Unmarshal([]byte(mainstr), &(data.([]interface{})[0]))
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	} else {
		data = (value)
	}
	return data
}

func init() {
	yiiMiddleman := new(YiiMiddleman)
	manager.MiddlemanRegister("yii", yiiMiddleman)
}
