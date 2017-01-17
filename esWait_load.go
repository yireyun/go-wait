// esWait_load.go

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
