package retry

import (
	"time"

	"github.com/hqmin9527/kits-go/src/logger"
)

// 指定重试时间组
func Do(f func() bool, rules []int) (ret bool) {
	for i := 0; ; i++ {
		ret = f()
		log(ret, i)
		if ret {
			return
		}
		if i >= len(rules) {
			return
		}
		time.Sleep(time.Duration(rules[i]) * time.Second)
	}
}

// 固定间隔重试
func TickDo(f func() bool, tick int) {
	for i := 0; ; i++ {
		ret := f()
		log(ret, i)
		if ret {
			return
		}
		time.Sleep(time.Duration(tick) * time.Second)
	}
}

func WithTimeout(f func() bool, tick int, d time.Duration) (ret bool) {
	deadline := time.Now().Add(d)
	for i := 0; time.Now().Before(deadline); i++ {
		ret = f()
		log(ret, i)
		if ret {
			return
		}
		time.Sleep(time.Duration(tick) * time.Second)
	}
	return false
}

func log(ret bool, i int) {
	if i == 0 {
		return
	}
	if ret {
		logger.Info("retry success at %d", i)
	} else {
		logger.Warn("retry failed at %d", i)
	}
}
