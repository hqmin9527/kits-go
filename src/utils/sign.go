package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func Sign(params map[string]string, key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("sign key is empty")
	}
	// 过滤掉signature参数
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signature" || len(k) == 0 {
			continue
		}
		keys = append(keys, k)
	}
	// 排序
	sort.Strings(keys)
	// 拼接参数
	keyValues := make([]string, len(keys))
	for i, k := range keys {
		keyValues[i] = k + "=" + params[k]
	}
	var p = strings.Join(keyValues, "&")
	// 签名
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(p))
	var signature = fmt.Sprintf("%X", mac.Sum(nil))
	return signature, nil
}

func SignQuery(query url.Values, key string) (string, error) {
	params := make(map[string]string)
	for k, v := range query {
		params[k] = v[0]
	}
	return Sign(params, key)
}

// 签名：一般性的map
// map的value只支持简单类型（不允许有json化歧义）
func SignMap(params map[string]any, key string) (string, error) {
	data := make(map[string]string, len(params))
	for k, v := range params {
		var valStr string
		switch val := v.(type) {
		case string:
			valStr = val
		default:
			d, _ := json.Marshal(val)
			valStr = string(d)
		}
		data[k] = valStr
	}
	return Sign(data, key)
}
