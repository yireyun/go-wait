// esWait

// +build arm arm64 mips64 mips64le ppc64 ppc64le

package wait

import (
	"sync/atomic"
	"unsafe"
)

//获取等待管道
func (w *EsWait) Wait() <-chan struct{} {
	ws := (*waitSignal)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws))))
	return *ws
}

//发送信号
func (w *EsWait) Signal() {
	var nws waitSignal = make(chan struct{})
	ws := (*waitSignal)(atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws)), unsafe.Pointer(&nws)))
	close(*ws) //go15，不关闭,性能更好
}
