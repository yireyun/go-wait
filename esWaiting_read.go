// esWaiting_read.go

// +build 386 amd64 amd64p32 s390x

package wait

import (
	"sync/atomic"
)

//获取等待管道
func (w *EsWaiting) Wait() <-chan struct{} {
	ws := w.ws
	atomic.AddUintptr(&ws.waiting, 1)
	return ws.c
}

//获取等待管道数量
func (w *EsWaiting) Waiting() int {
	ws := w.ws
	return int(atomic.LoadUintptr(&ws.waiting))
}

