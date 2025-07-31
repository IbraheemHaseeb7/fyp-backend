package jobs

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/IbraheemHaseeb7/fyp-backend/utils"
)

// Task represents a unit of work for the background queue.
type Task struct {
	Request utils.InternalApiRequest
}

// WorkerPool manages a pool of worker goroutines to process tasks.
type WorkerPool struct {
	Queue           []Task
	Size            int
	wg              sync.WaitGroup
	mutex           sync.Mutex // Protects Queue and isStopped
	taskChan        chan Task // Use this to signal new tasks
	numWorkers      int32     // Total number of worker goroutines (constant)
	numWorkersBusy  int32     // Number of workers currently processing a task
	workerPoolStats *sync.Map
}

// NewWorkerPool creates a new worker pool with the given size.
func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{
		Queue:           make([]Task, 0),
		Size:            size,
		taskChan:        make(chan Task), // Initialize the channel
		numWorkers:      int32(size), // Set the total number of workers
		numWorkersBusy:  0,
		workerPoolStats: &sync.Map{},
	}
	wp.workerPoolStats.Store("QueueLength", 0)
	wp.workerPoolStats.Store("BusyWorkers", 0)
	return wp
}

// Start starts the worker pool, launching the worker goroutines.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.Size; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop gracefully stops the worker pool by closing the queue and waiting for workers to finish.
func (wp *WorkerPool) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()
}

// EnqueueTask adds a new task to the queue or sends it directly to a worker.
func (wp *WorkerPool) EnqueueTask(task Task) {
	if atomic.LoadInt32(&wp.numWorkersBusy) < wp.numWorkers {

		select {
		case wp.taskChan <- task:
			atomic.AddInt32(&wp.numWorkersBusy, 1)
			wp.workerPoolStats.Store("BusyWorkers", atomic.LoadInt32(&wp.numWorkersBusy))
			return 

		default:
		}
	}

	// Queue the task if no worker was immediately available.
	wp.mutex.Lock()
	wp.Queue = append(wp.Queue, task) // Add to the queue
	wp.workerPoolStats.Store("QueueLength", len(wp.Queue))
	wp.mutex.Unlock()
}

// worker is the function executed by each worker goroutine.
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		task, ok := <-wp.taskChan // Receive tasks from the channel
		if !ok {
			fmt.Println("Worker stopped")
			return
		}
		wp.processTask(task)
	}
}

// processTask performs the internal API call for a given task.
func (wp *WorkerPool) processTask(task Task) {
	task.Request.InternalCall()
	atomic.AddInt32(&wp.numWorkersBusy, -1) // Decrement when finished
	wp.workerPoolStats.Store("BusyWorkers", atomic.LoadInt32(&wp.numWorkersBusy))
	wp.mutex.Lock()
	if len(wp.Queue) > 0 {
		//try to send another task.
		select{
		case wp.taskChan <- wp.Queue[0]:
			wp.Queue = wp.Queue[1:]
			wp.workerPoolStats.Store("QueueLength", len(wp.Queue))
			wp.mutex.Unlock()
			atomic.AddInt32(&wp.numWorkersBusy, 1)
			wp.workerPoolStats.Store("BusyWorkers", atomic.LoadInt32(&wp.numWorkersBusy))
		default:
			wp.mutex.Unlock()

		}
	} else {
		wp.mutex.Unlock()
	}
}

// NewTask creates a new task with the given request.
func NewTask(request utils.InternalApiRequest) Task {
	return Task{
		Request: request,
	}
}

// GetWorkerPoolStats returns the number of total and busy workers.
func (wp *WorkerPool) GetWorkerPoolStats() (int32, int32) {
	return atomic.LoadInt32(&wp.numWorkers), atomic.LoadInt32(&wp.numWorkersBusy)
}

func (wp *WorkerPool) GetQueueLength() int {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	return len(wp.Queue)
}
