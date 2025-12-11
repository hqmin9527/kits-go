package grace_stop

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hqmin9527/kits-go/src/logger"
)

// 监听关闭信号（pm2发送SIGINT，k8s发送SIGTERM），然后关闭stop协程
// 程序可以调用Closed来获取stop的状态
// 也可以通过GetStopChan来拿到stop，自己监听
// Once.Do保证MonitorStop只会被调用一次

var stop = make(chan struct{}) // 优雅关闭程序标志
var once = sync.Once{}

func MonitorStop() {
	once.Do(func() {
		// 收到pm2或者k8s的Signal后，设置stop为true，处理完当前的历史课堂后，程序结束
		signals := make(chan os.Signal, 1)

		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		sig := <-signals
		logger.Info("process will stop gracefully for sig[%d]", sig)
		close(stop)
	})
}

func Closed() bool {
	select {
	case <-stop:
		return true
	default:
		return false
	}
}

func GetStopChan() <-chan struct{} {
	return stop
}
