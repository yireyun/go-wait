// esWait_test.go
package wait

import (
	"sync"
	"testing"
	"time"
)

func TestWaiting(t *testing.T) {
	w := NewWaiting()
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
			if w.Signal() > 0 {
				t.Logf("Signal  No:%d\n", no)
			}
		}(i)
	}
	wg.Wait()
}

var waiting = NewWaiting()

func WaitingWaitSignalAsync(N, waitGrp, signalGrp, Gnt int, sleep bool, logf func(format string, args ...interface{})) {
	var wg1 sync.WaitGroup
	stop := false
	wg1.Add(waitGrp)
	waitCnt := make([]int, waitGrp)
	start := time.Now()
	for g := 0; g < waitGrp; g++ {
		go func(cnt *int) {
			defer wg1.Done()
			for i := 0; i < Gnt; i++ {
				waiting.Wait()
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
				if n := waiting.Signal(); n > 0 {
					*cnt++
				}
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

func WaitingWaitSignalSync(N, waitGrp, signalGrp, Gnt int, sleep bool, logf func(format string, args ...interface{})) {
	var wg1 sync.WaitGroup
	stop := false
	wg1.Add(waitGrp)
	waitCnt := make([]int, waitGrp)
	start := time.Now()
	for g := 0; g < waitGrp; g++ {
		go func(cnt *int) {
			defer wg1.Done()
			for i := 0; i < Gnt; i++ {
				<-waiting.Wait()
				*cnt++
			}
		}(&waitCnt[g])
	}
	var wg2 sync.WaitGroup
	wg2.Add(signalGrp)
	signalCnt := make([]int, signalGrp)
	waitingCnt := make([]int, waitGrp)
	if sleep {
		time.Sleep(time.Microsecond)
	}
	for g := 0; g < signalGrp; g++ {
		go func(cnt *int) {
			defer wg2.Done()
			for !stop {
				if n := waiting.Signal(); n > 0 {
					*cnt++
					waitingCnt[n-1]++
				}
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
	waiting := 0
	for i := range waitingCnt {
		waiting += waitingCnt[i]
	}
	waitOp := use / time.Duration(wait)
	logf("Times:%10d, WaitGrp:%3d, SignalGrp:%3d, Gnt:%10d, Use:%5dms %7dns/op, Wait:%10d %7dns/op, Signal:%10d %7dns/op, Waiting:%10d%+6v\n",
		N, waitGrp, signalGrp, Gnt, use/time.Millisecond, op, wait, waitOp, signal, signalOp, waiting, waitingCnt)
}

func WaitingWaitSignalSyncAll(N, waitGrp, signalGrp, Gnt int, sleep bool, logf func(format string, args ...interface{})) {
	var wg1 sync.WaitGroup
	stop := false
	wg1.Add(waitGrp)
	waitCnt := make([]int, waitGrp)
	start := time.Now()
	for g := 0; g < waitGrp; g++ {
		go func(cnt *int) {
			defer wg1.Done()
			for i := 0; i < Gnt; i++ {
				<-waiting.Wait()
				*cnt++
			}
		}(&waitCnt[g])
	}
	var wg2 sync.WaitGroup
	wg2.Add(signalGrp)
	signalCnt := make([]int, signalGrp)
	waitingCnt := make([]int, waitGrp)
	if sleep {
		time.Sleep(time.Microsecond)
	}
	for g := 0; g < signalGrp; g++ {
		go func(cnt *int) {
			defer wg2.Done()
			for !stop {
				if w := waiting.Waiting(); w == waitGrp {
					if n := waiting.Signal(); n > 0 {
						*cnt++
						waitingCnt[n-1]++
					}
				}
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
	waiting := 0
	for i := range waitingCnt {
		waiting += waitingCnt[i]
	}
	waitOp := use / time.Duration(wait)
	logf("Times:%10d, WaitGrp:%3d, SignalGrp:%3d, Gnt:%10d, Use:%5dms %7dns/op, Wait:%10d %7dns/op, Signal:%10d %7dns/op, Waiting:%10d%+6v\n",
		N, waitGrp, signalGrp, Gnt, use/time.Millisecond, op, wait, waitOp, signal, signalOp, waiting, waitingCnt)
}

func BenchmarkWaitingWaitDo(b *testing.B) {
	for i := 1; i < b.N; i++ {
		waiting.Wait()
	}
}

func BenchmarkWaitingWaitGo(b *testing.B) {
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
				waiting.Wait()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkWaitingSignalDo(b *testing.B) {
	for i := 1; i < b.N; i++ {
		waiting.Signal()
	}
}

func BenchmarkWaitingSignalGo(b *testing.B) {
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
				waiting.Signal()
			}
		}()
	}
	wg.Wait()
}

func TestWaitingWaitSignalSync(t *testing.T) {
	N := 1000
	for s := 1; s <= 1; s = s * 2 {
		for w := 1; w <= 128; w = w * 2 {
			WaitingWaitSignalSync(N, w, s, N, false, t.Logf)
		}
	}
}

func TestWaitingWaitSignalSyncAll(t *testing.T) {
	N := 100
	for s := 1; s <= 1; s = s * 2 {
		for w := 1; w <= 128; w = w * 2 {
			WaitingWaitSignalSyncAll(N, w, s, N, false, t.Logf)
		}
	}
}

func TestWaitingWaitSignalAsync(t *testing.T) {
	N := 10000 * 10
	for s := 1; s <= 1; s = s * 2 {
		for w := 1; w <= 128; w = w * 2 {
			WaitingWaitSignalAsync(N, w, s, N, false, t.Logf)
		}
	}
}

func BenchmarkWaitingWaitSignalSync(b *testing.B) {
	Grp := 1
	Gnt := b.N
	WaitingWaitSignalSync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitingWaitSignalOneSync(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}

	WaitingWaitSignalSync(b.N, Grp, 1, Gnt, false, b.Logf)
}

func BenchmarkWaitingWaitSignalMulSync(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}

	WaitingWaitSignalSync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitingWaitSignalAsync(b *testing.B) {
	Grp := 1
	Gnt := b.N
	WaitingWaitSignalAsync(b.N, Grp, Grp, Gnt, false, b.Logf)
}

func BenchmarkWaitingWaitSignalOneAsync(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}

	WaitingWaitSignalAsync(b.N, Grp, 1, Gnt, false, b.Logf)
}

func BenchmarkWaitingWaitSignalMulAsync(b *testing.B) {
	Grp := 100
	Gnt := b.N / Grp
	if Gnt == 0 {
		Gnt = b.N
		Grp = 1
	}

	WaitingWaitSignalAsync(b.N, Grp, Grp, Gnt, false, b.Logf)
}
