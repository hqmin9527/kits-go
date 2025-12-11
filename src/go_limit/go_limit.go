package go_limit

import (
	"runtime/debug"
	"sync"

	"github.com/hqmin9527/kits-go/src/logger"
)

type syncErrList struct {
	errs []error
	mu   sync.RWMutex
}

func (l *syncErrList) appendError(err error) {
	l.mu.Lock()
	l.errs = append(l.errs, err)
	l.mu.Unlock()
}

func (l *syncErrList) listErrors() []error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.errs
}

// GoLimit 限制一个同步任务中并发的协程数
type GoLimit struct {
	c    chan struct{}
	wg   *sync.WaitGroup
	errs *syncErrList
}

func New(size int) *GoLimit {
	return &GoLimit{
		c:    make(chan struct{}, size),
		wg:   &sync.WaitGroup{},
		errs: &syncErrList{},
	}
}

func (g *GoLimit) Run(f func()) *GoLimit {
	g.wg.Add(1)
	g.c <- struct{}{}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("GoLimit Run panic as: %s\n%s", err, debug.Stack())
			}
			g.wg.Done()
			<-g.c
		}()
		f()
	}()
	return g
}

func (g *GoLimit) RunError(f func() error) *GoLimit {
	g.wg.Add(1)
	g.c <- struct{}{}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("GoLimit RunError panic as: %s\n%s", err, debug.Stack())
			}
			g.wg.Done()
			<-g.c
		}()
		if err := f(); err != nil {
			g.errs.appendError(err)
		}
	}()
	return g
}

func (g *GoLimit) Wait() {
	g.wg.Wait()
}

func (g *GoLimit) ListErrors() []error {
	return g.errs.listErrors()
}

func (g *GoLimit) FirstError() error {
	errs := g.errs.listErrors()
	if len(errs) > 0 {
		return errs[0]
	} else {
		return nil
	}
}
