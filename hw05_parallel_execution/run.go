package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	taskCh := make(chan Task, n)
	var errorCount int32

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if int(atomic.LoadInt32(&errorCount)) >= m {
					return
				}
				res := task()
				// resCh <- res
				if res != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(tasks); i++ {
			if int(atomic.LoadInt32(&errorCount)) >= m {
				break
			}
			taskCh <- tasks[i]
		}
		close(taskCh)
	}()

	wg.Wait()
	if int(atomic.LoadInt32(&errorCount)) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
