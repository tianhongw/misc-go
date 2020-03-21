package util

import (
	"errors"
	"runtime"
	"sync"

	"github.com/tianhongw/misc-go/log"
	"go.uber.org/zap"
)

// WorkerPool is pool of workers
type WorkerPool struct {
	wgWrapper *WaitGroupWrapper
	doneChan  chan struct{}
	workChan  chan func()
	once      sync.Once
}

func NewWorkerPool() *WorkerPool {
	wp := &WorkerPool{doneChan: make(chan struct{}), workChan: make(chan func())}
	wp.Start()
	return wp
}

// Start worker pool
func (wp *WorkerPool) Start() {
	n := runtime.NumCPU()
	for i := 0; i < n; i++ {
		wp.wgWrapper.Wrap(func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Instance().Error("workerPool", zap.Any("err", err))
				}
			}()
			for {
				select {
				case f := <-wp.workChan:
					f()
				case <-wp.doneChan:
					return
				}
			}
		})
	}
}

// Close the worker pool
func (wp *WorkerPool) Close() {
	wp.once.Do(func() {
		close(wp.doneChan)
		wp.wgWrapper.Wait()
	})
}

var (
	// ErrWorkerPoolClosed when run on closed pool
	ErrWorkerPoolClosed = errors.New("workerPool closed")
)

// Run a task
func (wp *WorkerPool) Run(f func()) error {
	select {
	case <-wp.doneChan:
		return ErrWorkerPoolClosed
	default:
	}

	select {
	case wp.workChan <- f:
		return nil
	case <-wp.doneChan:
		return ErrWorkerPoolClosed
	}
}
