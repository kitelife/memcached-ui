package middleman

import (
	"github.com/picasso250/memcached-ui/middleman/manager"
)

type DefaultMiddleman struct{}

func (dmm *DefaultMiddleman) Config(config map[string]string) bool {
	return true
}

func (dmm DefaultMiddleman) GenInnerKey(key string) string {
	return key
}

func (dmm DefaultMiddleman) SerializeValue(value string) string {
	return value
}

func (dmm DefaultMiddleman) UnserializeValue(value string) interface{} {
	return value
}

func init() {
	defaultMiddleman := new(DefaultMiddleman)
	manager.MiddlemanRegister("default", defaultMiddleman)
}
