// esWait_read.go

// +build 386 amd64 amd64p32 s390x

package wait

//获取等待管道
func (w *EsWait) Wait() <-chan struct{} {
	return *w.ws
}
