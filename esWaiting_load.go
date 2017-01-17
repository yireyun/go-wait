// esWait

// +build arm arm64 mips64 mips64le ppc64 ppc64le

package wait

import (
	"sync/atomic"
	"unsafe"
)

//获取等待管道
func (w *EsWaiting) Wait() <-chan struct{} {
	ws := (*waitingSignal)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws))))
	atomic.AddUintptr(&ws.waiting, 1)
	return ws.c
}

//获取等待管道数量
func (w *EsWaiting) Waiting() int {
	ws := (*waitingSignal)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws))))
	return int(atomic.LoadUintptr(&ws.waiting))
}

//发送信号
func (w *EsWaiting) Signal() int {
	fs := waitingPool.Get().(*waitingSignal)
	ws := (*waitingSignal)(atomic.SwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&w.ws)), unsafe.Pointer(fs)))
	if waiting := atomic.LoadUintptr(&ws.waiting); waiting > 0 {
		close(ws.c) //go15，不关闭,性能更好
		return int(waiting)
	} else {
		waitingPool.Put(ws)
		return 0
	}
}
