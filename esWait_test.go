// esWait_test.go
package wait

import (
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	w := NewWait()
	N := 100
	var wg sync.WaitGroup

	for i := 1; i <= N; i++ {
		wg.Add(1)
		go func(no int) {
			defer wg.Done()
			start := time.Now()
			end := start
			c := w.Wait()
			select {
			case _, ok := <-c:
				end = time.Now()
				if !ok {
					t.Logf("Waiting No:%03d, Chan:%v, Use:%v\n", no, c, end.Sub(start))
				}
			case <-time.After(5 * time.Microsecond):
				end = time.Now()
				t.Logf("Timeout No:%03d, Chan:%v, Use:%v\n", no, c, end.Sub(start))
			}

		}(i)
	}
	for i := 1; i <= N; i++ {
		wg.Add(1)
		go func(no int) {
			defer wg.Done()
			w.Signal()
			t.Logf("Signal  No:%d\n", no)
		}(i)
	}
	wg.Wait()
}

var wait = NewWait()

func WaitSignalAsync(N, waitGrp, signalGrp, Gnt int, sleep bool, logf func(format string, args ...interface{})) {
	var wg1 sync.WaitGroup
	stop := false
	wg1.Add(waitGrp)
	waitCnt := make([]int, waitGrp)
	start := time.Now()
	for g := 0; g < waitGrp; g++ {
		go func(cnt *int) {
			defer wg1.Done()
			for i := 0; i < Gnt; i++ {
				wait.Wait()
				*cnt++
			}
		}(&waitCnt[g])
	}
	var wg2 sync.WaitGroup
	wg2.Add(signalGrp)
	signalCnt := make([]int, signalGrp)
	if sleep {
		time.Sleep(time.Microsecond)
	}
	for g := 0; g < signalGrp; g++ {
		go func(cnt *int) {
			defer wg2.Done()
			for !stop {
				wait.Signal()
				*cnt++
			}
		}(&signalCnt[g])
	}
	wg1.Wait()
	stop = true
	wg2.Wait()
	end := time.Now()
	use := end.Sub(start)
	op := use / time.Duration(N)
	signal := 0
	for i := range signalCnt {
		signal += signalCnt[i]
	}
	signalOp := use / time.Duration(signal)
	wait := 0
	for i := range waitCnt {
		wait += waitCnt[i]
	}
	waitOp := use / time.Duration(wait)
	logf("Times:%10d, WaitGrp:%3d, SignalGrp:%3d, Gnt:%10d, Use:%5dms %7dns/op, Wait:%10d %7dns/op, Signal:%10d %7dns/op\n",
		N, waitGrp, signalGrp, Gnt, use/time.Millisecond, op, wait, waitOp, signal, signalOp)
}

func WaitSignalSync(N, waitGrp, signalGrp, Gnt int, sleep bool, logf func(format string, args ...interface{})) {
	var wg1 sync.WaitGroup
	stop := false
	wg1.Add(waitGrp)
	waitCnt := make([]int, waitGrp)
	start := time.Now()
	for g := 0; g < waitGrp; g++ {
		go func(cnt *int) {
			defer wg1.Done()
			for i := 0; i < Gnt; i++ {
				<-wait.Wait()
				*cnt++
			}
		}(&waitCnt[g])
	}
	var wg2 sync.WaitGroup
	wg2.Add(signalGrp)
	signalCnt := make([]int, signalGrp)
	if sleep {
		time.Sleep(time.Microsecond)
	}
	for g := 0; g < signalGrp; g++ {
		go func(cnt *int) {
			defer wg2.Done()
			for !stop {
				wait.Signal()
				*cnt++
			}
		}(&signalCnt[g])
	}
	wg1.Wait()
	stop = true
	wg2.Wait()
	end := time.Now()
	use := end.Sub(start)
	op := use / time.Duration(N)
	signal := 0
	for i := range signalCnt {
		signal += signalCnt[i]
	}
	signalOp := use / time.Duration(signal)
	wait := 0
	for i := range waitCnt {
		wait += waitCnt[i]
	}
	waitOp := use / time.Duration(wait)
	logf("Times:%10d, WaitGrp:%3d, SignalGrp:%3d, Gnt:%10d, Use:%5dms %7dns/op, Wait:%10d %7dns/op, Signal:%10d %7dns/op\n",
		N, waitGrp, signalGrp, Gnt, use/time.Millisecond, op, wait, waitOp, signal, signalOp)
}
func BenchmarkWaitDo(b *testing.B) {
	for i := 1; i < b.N; i++ {
		wait.Wait()
	}
}

func BenchmarkWaitGo(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}
	var wg sync.WaitGroup
	wg.Add(Grp)
	for g := 0; g < Grp; g++ {
		go func() {
			defer wg.Done()
			for i := 1; i < Gnt; i++ {
				wait.Wait()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkSignalDo(b *testing.B) {
	for i := 1; i < b.N; i++ {
		wait.Signal()
	}
}

func BenchmarkSignalGo(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}
	var wg sync.WaitGroup
	wg.Add(Grp)
	for g := 0; g < Grp; g++ {
		go func() {
			defer wg.Done()
			for i := 1; i < Gnt; i++ {
				wait.Signal()
			}
		}()
	}
	wg.Wait()
}

func TestWaitSignalSync(t *testing.T) {
	N := 1000
	for s := 1; s <= 1; s = s * 2 {
		for w := 1; w <= 128; w = w * 2 {
			WaitSignalSync(N, w, s, N, false, t.Logf)
		}
	}
}

func TestWaitSignalAsync(t *testing.T) {
	N := 10000 * 10
	for s := 1; s <= 1; s = s * 2 {
		for w := 1; w <= 128; w = w * 2 {
			WaitSignalAsync(N, w, s, N, false, t.Logf)
		}
	}
}
func BenchmarkWaitSignalSyncDo(b *testing.B) {
	Grp := 1
	Gnt := b.N
	WaitSignalSync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitSignalSyncGo(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}
	WaitSignalSync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitSignalAsyncDo(b *testing.B) {
	Grp := 1
	Gnt := b.N
	WaitSignalAsync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitSignalAsyncGo(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}
	WaitSignalAsync(b.N, Grp, Grp, Gnt, false, b.Logf)
}
