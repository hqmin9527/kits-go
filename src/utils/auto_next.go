package utils

import "sync"

func AutoNext(start int) (func() int, func(int)) {
	next := func() int {
		res := start
		start++
		return res
	}
	setNext := func(next int) {
		start = next
	}
	return next, setNext
}

func ConcurrentAutoNext(start int) func() int {
	l := sync.Mutex{}
	return func() int {
		l.Lock()
		res := start
		start++
		l.Unlock()
		return res
	}
}
