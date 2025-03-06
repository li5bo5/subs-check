package utils

import (
	"sync"
)

// Task 是一个任务结构体，包含任务参数和结果
type Task struct {
	ID     int
	Args   interface{}
	Result interface{}
	Err    error
}

// ThreadPool 是线程池结构体
type ThreadPool struct {
	NumWorkers int
	Tasks      chan Task
	Results    chan Task
	Wg         sync.WaitGroup
	TaskFunc   func(interface{}) (interface{}, error)
	results    []Task
	mu         sync.Mutex
	resultWg   sync.WaitGroup
}

// NewThreadPool 创建一个新的线程池
func NewThreadPool(numWorkers int, taskFunc func(interface{}) (interface{}, error)) *ThreadPool {
	tp := &ThreadPool{
		NumWorkers: numWorkers,
		Tasks:      make(chan Task, numWorkers),
		Results:    make(chan Task, numWorkers),
		TaskFunc:   taskFunc,
		results:    make([]Task, 0),
	}
	go tp.processResults()
	return tp
}

// Start 启动线程池中的工作线程
func (tp *ThreadPool) Start() {
	for i := 0; i < tp.NumWorkers; i++ {
		tp.Wg.Add(1)
		go tp.worker()
	}
}

// worker 是工作线程，从Tasks通道中获取任务并执行
func (tp *ThreadPool) worker() {
	defer tp.Wg.Done()
	for task := range tp.Tasks {
		task.Result, task.Err = tp.TaskFunc(task.Args)
		tp.Results <- task
	}
}

// processResults 消费 Results 通道中的结果，并将其存储到 results 切片中
func (tp *ThreadPool) processResults() {
	for result := range tp.Results {
		tp.mu.Lock()
		tp.results = append(tp.results, result)
		tp.mu.Unlock()
		tp.resultWg.Done()
	}
}

// GetResults 获取所有任务的结果
func (tp *ThreadPool) GetResults() []Task {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.results
}

// AddTaskArgs 根据参数列表创建任务并批量添加到线程池
func (tp *ThreadPool) AddTaskArgs(argsList []interface{}) {
	for i, args := range argsList {
		task := Task{
			ID:   i + 1,
			Args: args,
		}
		tp.resultWg.Add(1)
		tp.Tasks <- task
	}
}

// Wait 等待所有任务完成并关闭通道
func (tp *ThreadPool) Wait() {
	close(tp.Tasks)
	tp.Wg.Wait()
	tp.resultWg.Wait()
	close(tp.Results)
}
