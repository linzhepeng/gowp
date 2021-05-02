package gowp

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RUNNING = 1
	STOPPED = 0
)

type Job struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

type Pool struct {
	Cap            uint64
	RunningWorkers uint64
	Status         int64
	JobChan        chan *Job
	PanicHandler   func(v ...interface{})
	Lock           sync.Mutex
}

// NewPool 初始化协程池
func NewPool(cap uint64) (*Pool, error) {
	if cap <= 0 {
		return nil, errors.New("cap must be greater than zero. ")
	}
	pool := &Pool{
		Cap:     cap,
		Status:  RUNNING,
		JobChan: make(chan *Job, cap),
	}

	return pool, nil
}

// GetRunningWorkers 原子操作，为了协程安全
func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.RunningWorkers)
}

func (p *Pool) AddRunningWorkers() {
	atomic.AddUint64(&p.RunningWorkers, 1)
}

func (p *Pool) DelRunningWorkers() {
	atomic.AddUint64(&p.RunningWorkers, ^uint64(0))
}

// Put 将任务信息加入任务通道中
func (p *Pool) Put(job *Job) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if p.Status == STOPPED {
		return errors.New("pool is already closed")
	}

	if p.Status == RUNNING {
		p.JobChan <- job
	}

	if p.GetRunningWorkers() < p.Cap {
		p.Run()
	}
	return nil
}

// Run 处理任务
func (p *Pool) Run() {
	p.AddRunningWorkers()
	go func() {
		defer func() {
			p.DelRunningWorkers()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("Worker panic: %s\n", r)
				}
			}
		}()

		for {
			select {
			case job, ok := <-p.JobChan:
				if !ok {
					return
				}
				job.Handler(job.Params...)
			}
		}
	}()
}

func (p *Pool) SetStatus(status int64) bool {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if p.Status == status {
		return false
	}

	p.Status = status
	return true
}

func (p *Pool) ClosePool() {
	if !p.SetStatus(STOPPED) {
		return
	}
	for len(p.JobChan) > 0 {
		time.Sleep(time.Second)
	}

	close(p.JobChan)
}
