package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// 分布式id

func GenObjectId() string {
	return bson.NewObjectId().Hex()
}

// 默认生成v7的uuid
func GenUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

// 解析ObjectID中的时间戳(ms)
func ParseObjectIDTs(objectId string) (int64, error) {
	if len(objectId) != 24 {
		return 0, errors.New("invalid ObjectId length")
	}
	dataStr := objectId[:8]
	result, err := strconv.ParseInt(dataStr, 16, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse int")
	}
	return result * 1000, nil
}

// 解析uuid中的时间戳(ms)
func ParseUUIDTs(uuid string) (int64, error) {
	clean := strings.ReplaceAll(uuid, "-", "")
	if len(clean) != 32 {
		return 0, errors.New("invalid UUID length")
	}

	version := clean[12]

	switch version {
	case '7':
		return parseV7(clean)
	default:
		return 0, errors.New("unsupported UUID version")
	}
}

func parseV7(clean string) (int64, error) {
	tsHex := clean[:12]
	ms, err := strconv.ParseInt(tsHex, 16, 64)
	if err != nil {
		return 0, err
	}
	return ms, nil
}
