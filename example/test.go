package main

import (
	"fmt"
	"github.com/linzhepeng/gowp"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// 创建容量为 10 的任务池
	pool, err := gowp.NewPool(10)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		// 创建任务
		job := &gowp.Job{
			Handler: func(v ...interface{}) {
				wg.Done()
				fmt.Println(v)
			},
			Params: []interface{}{i, i * 2, "hello"},
		}
		// 将任务放入任务池
		_ = pool.Put(job)
	}
	wg.Wait()
	// 安全关闭任务池（保证已加入池中的任务被消费完）
	pool.ClosePool()
}
