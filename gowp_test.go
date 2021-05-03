package gowp

import (
	"sync"
	"sync/atomic"
	"testing"
)

var sum int64
var runTimes = 1000000
var wg = sync.WaitGroup{}
var lock sync.Mutex

func demoJob(v ...interface{}) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		lock.Lock()
		sum += int64(1)
		lock.Unlock()
	}
}

func demoJob2(v ...interface{}) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

//func BenchmarkGoroutineMutex(b *testing.B) {
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		go demoJob()
//	}
//	wg.Wait()
//}
//
//func BenchmarkPoolMutex(b *testing.B) {
//	pool, err := NewPool(20)
//	if err != nil {
//		b.Error(err)
//	}
//
//	job := &Job{
//		Handler: demoJob,
//	}
//
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		pool.Put(job)
//	}
//	wg.Wait()
//}

func BenchmarkGoroutineAtomic(b *testing.B) {
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go demoJob2()
	}
	wg.Wait()
}

func BenchmarkPoolAtomic(b *testing.B) {
	pool, err := NewPool(20)
	if err != nil {
		b.Error(err)
	}

	job := &Job{
		Handler: demoJob2,
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		pool.Put(job)
	}

	wg.Wait()
}
