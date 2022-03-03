package cluster

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	dec = iota - 1
	inc
	Running = true
	Stopped = false
)

// 协程池结构
type pool struct {
	sync.Mutex
	capacity int64      // 协程池容量
	running  int64      // 协程池装载量
	status   bool       // 状态
	taskCh   chan *task // 任务管道
}

func NewPool(capacity int64) (*pool, error) {
	if capacity > 0 {
		return &pool{
			capacity: capacity,
			status:   Running,
			taskCh:   make(chan *task, capacity),
		}, nil
	}
	return nil, errors.New("capacity must be more than 0")
}

// PutTask 协程池添加任务
func (p *pool) PutTask(task *task) error {
	if !p.status { // 检查状态
		return errors.New("task pool is stopped")
	}
	p.Lock()
	defer func() {
		p.Unlock()
	}()

	running := p.getRunning()
	if running < p.running { //还有空余容量创建执行工作任务
		p.run()
	} else {
		p.taskCh <- task
	}

	return nil
}

// Close 协程池关闭方法
func (p *pool) Close() {
	p.status = Stopped
	for len(p.taskCh) > 0 {
	}
	close(p.taskCh)
}

//
func (p *pool) run() {
	p.alertRunning(inc)
	go func() {
		defer func() {
			p.alertRunning(dec)
		}()

		for {
			select {
			case task, ok := <-p.taskCh:
				if ok {
					task.Handle(task.Params)
				}
			}
		}

	}()
}

// 修改协程池运行数量返回当前已有协程数
func (p *pool) alertRunning(value int64) int64 {
	return atomic.AddInt64(&(p.running), value)
}

// 获取当前以后协程数
func (p *pool) getRunning() int64 {
	return atomic.LoadInt64(&(p.running))
}
