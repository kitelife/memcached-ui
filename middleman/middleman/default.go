package middleman

import (
	"github.com/youngsterxyf/memcached-ui/middleman/manager"
)

type DefaultMiddleman struct{}

func (dmm DefaultMiddleman) GenInnerKey(key string, config interface{}) string {
	return key
}

func (dmm DefaultMiddleman) SerializeValue(value string, config interface{}) string {
	return value
}

func (dmm DefaultMiddleman) UnserializeValue(value string, config interface{}) string {
	return value
}

func init() {
	defaultMiddleman := DefaultMiddleman{}
	manager.MiddlemanRegister("default", defaultMiddleman)
}
