// esWait

// +build 386 amd64 amd64p32 s390x

package wait

import (
	"sync/atomic"
	"unsafe"
)

//获取等待管道
func (w *EsWait) Wait() <-chan struct{} {
	return *w.ws
}

//发送信号
func (w *EsWait) Signal() {
	var nws waitSignal = make(chan struct{})
	ws := (*waitSignal)(atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws)), unsafe.Pointer(&nws)))
	close(*ws) //go15，不关闭,性能更好
}
