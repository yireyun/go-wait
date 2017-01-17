// esWait
package wait

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
