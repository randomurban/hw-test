package hw05parallelexecution

import (
	"errors"
	"sync"
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
	errorCount := 0
	goodCount := 0

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if errorCount >= m {
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
				errorCount++
				if errorCount >= m {
					return
				}
			} else {
				goodCount++
			}
			if goodCount+errorCount >= len(tasks) {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(tasks); i++ {
			if errorCount >= m {
				break
			}
			taskCh <- tasks[i]
		}
		close(taskCh)
	}()

	wg.Wait()
	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
