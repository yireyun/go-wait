// esWait
package wait

import (
	"sync"
)

type waitingSignal struct {
	c       chan struct{}
	waiting uintptr
}

//使用chan和atomic实现的一个等待连接可用
type EsWaiting struct {
	ws *waitingSignal //指向等待signal
}

var waitingPool = sync.Pool{New: func() interface{} {
	return &waitingSignal{c: make(chan struct{})}
}}

func NewWaiting() *EsWaiting {
	w := new(EsWaiting)
	w.ws = &waitingSignal{c: make(chan struct{})}
	return w
}
