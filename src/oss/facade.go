package oss

import "github.com/hqmin9527/kits-go/src/logger"

func NewWrapper(c *Config) (*Wrapper, error) {
	return newOssWrapper(c)
}

func GetWrapper(ossName string) *Wrapper {
	return _m[ossName]
}

// 仅在初始化时设置，没有并发问题
var _m = make(map[string]*Wrapper)

func InitByMap(ossMap map[string]*Config) {
	for k, v := range ossMap {
		logger.Debug("will init oss, name: %s, config: %s", k, v)
		w, err := NewWrapper(v)
		if err != nil {
			panic(err)
		}
		_m[k] = w
	}
}
