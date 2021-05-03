# gowp
- 高性能go协程池，代码清晰易懂
# 安装
`go get github.com/linzhepeng/gowp`
# 为什么要用协程池
- go的协程虽然轻量级，单个goroutine只有2KB，但毫无止境地使用也会占据大量的系统资源，导致调度性能下降、GC 频繁、内存暴涨, 引发一系列问题。
- 可以测试如下代码，无限开启goroutine很快会造成系统的崩溃
```
func main() {
    for i:=0;;i++ {
        go func(i int) {
            fmt.Println(i)
            time.Sleep(time.Second)
        }(i)
    }
}
```
- gowp限制了最多可启动的goroutine数量，并保持原有性能，在海量并发的场景中性能优势明显
# 原理
- 生产者消费者模型：
    
    1、创建一个容量为n的线程池（gpwp）
  
    2、生成任务（job）放入任务通道（jobChan）
  
    3、协程池未满时，每放入一个任务起一个goroutine，否则任务会被已有的n个goroutine抢占执行
# 使用方法
- 例子
```
package main

import (
"fmt"
"github.com/linzhepeng/gowp"
"sync"
)

var wg sync.WaitGroup

func main() {
    //初始化协程池，容量为10
    pool, err := gowp.NewPool(10)
    if err != nil {
        panic(err)
    }
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        //生成1000个任务
        job := &gowp.Job{
            Handler: func(v ...interface{}) {
                wg.Done()
                fmt.Println(v)
            },
            Params: []interface{}{i,i*2},
        }
    // 将任务放入协程池
        _ = pool.Put(job)
    }
    wg.Wait()
    //安全关闭协程池
    pool.ClosePool()
}
```
# benchmark性能测试
测试环境

```
goos:windows
goarch:amd64
cpu:Intel(R) Core(TM) i7-10875H CPU @ 2.30GHz
```

测试代码详见gowp_test.go文件

- 100w次加互斥锁的自增操作

| 操作模式 | 操作时间消耗ns/op | 内存分配大小B/op | 内存分配次数 |
| :-----: | -------------: | -------------: | --------: |
| 未使用协程池| 11872802400 |  312577056|1175045|
| 开启容量为20的协程池|5018586500| 456|19|

- 100w次原子增量操作

| 操作模式 | 操作时间消耗ns/op | 内存分配大小B/op | 内存分配次数 |
| :-----: | -------------: | -------------: | --------: |
| 未使用协程池| 1678428700 |  6950496|16375|
| 开启容量为20的协程池|1613707000| 3848 |40|


