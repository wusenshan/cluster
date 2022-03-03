package cluster

// Task 任务interface
type Task interface {
	Handler(params ...interface{})
}

// Job 工作结构体
type Job struct {
	Task   Task
	Params []interface{}
}

// 任务结构
type task struct {
	Handle func(...interface{})
	Params []interface{}
}

func NewTask(handle func(...interface{}), params []interface{}) *task {
	return &task{Handle: handle, Params: params}
}
