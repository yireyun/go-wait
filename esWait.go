// esWait
package wait

import (
	"sync/atomic"
	"unsafe"
)

type waitSignal chan struct{}

//使用chan和atomic实现的一个等待连接可用
type EsWait struct {
	ws *waitSignal //指向等待signal
}

func NewWait() *EsWait {
	w := new(EsWait)
	var ws waitSignal = make(chan struct{})
	w.ws = &ws
	return w
}

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
