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

	wg := sync.WaitGroup{}
	taskCh := make(chan Task, n)
	resCh := make(chan error, n)
	var errorCount int32
	var goodCount int32

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if int(errorCount) >= m {
					return
				}
				res := task()
				resCh <- res
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for res := range resCh {
			if res != nil {
				atomic.AddInt32(&errorCount, 1)
				if int(errorCount) >= m {
					return
				}
			} else {
				goodCount++
			}
			if int(goodCount+errorCount) >= len(tasks) {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(tasks); i++ {
			if int(errorCount) >= m {
				break
			}
			taskCh <- tasks[i]
		}
		close(taskCh)
	}()

	wg.Wait()
	if int(errorCount) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
