package utils

import (
	"strconv"
	"unicode/utf8"
)

// 掩码字符串
func MaskString(key string) string {
	l := len(key)
	switch {
	case l == 0:
		return "****"
	case l < 10:
		return key[:1] + "****" + key[l-1:]
	default:
		return key[:3] + "****" + key[l-3:]
	}
}

// 限制字符串长度（精确）
func LimitStringPrecise(s string, maxCharsSize int) string {
	if utf8.RuneCountInString(s) > maxCharsSize {
		r := []rune(s)
		return string(r[:maxCharsSize])
	}
	return s
}

// 限制字符串长度（带omit）
func LimitStringOmit(s string, maxCharsSize int) string {
	omitCnt := utf8.RuneCountInString(s) - maxCharsSize
	if omitCnt > 0 {
		r := []rune(s)
		return string(r[:maxCharsSize]) + "... omit " + strconv.Itoa(omitCnt) + " runes"
	}
	return s
}
