package utils

import (
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// 分布式id

func GenObjectId() string {
	return bson.NewObjectId().Hex()
}

func GenUUID() string {
	return uuid.NewV4().String()
}

// 解析mongo 的_id值
// 前8位16进制的以秒计算的时间值，转化成ms
func GetObjectTime(objectId string) (int64, bool) {
	if len(objectId) < 8 {
		return 0, false
	}
	dataStr := objectId[:8]
	result, err := strconv.ParseInt(dataStr, 16, 64)
	if err != nil {
		return 0, false
	}
	return result * 1000, true
}
