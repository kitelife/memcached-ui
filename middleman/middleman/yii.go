package middleman

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"

	"github.com/youngsterxyf/memcached-ui/middleman/manager"
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
	return value
}

func init() {
	yiiMiddleman := new(YiiMiddleman)
	manager.MiddlemanRegister("yii", yiiMiddleman)
}
