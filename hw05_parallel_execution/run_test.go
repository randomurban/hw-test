package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		withErrors(t, 50, 10, 23)
	})

	t.Run("with Errors if m<=0", func(t *testing.T) {
		withErrors(t, 5, 10, 0)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		withoutErrors(t, 50, 5, 1)
	})

	t.Run("without Errors if n>len(tasks)", func(t *testing.T) {
		withoutErrors(t, 5, 10, 1)
	})
}

func withoutErrors(t *testing.T, tasksCount int, workersCount int, maxErrorsCount int) {
	t.Helper()
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32
	var sumTime time.Duration

	for i := 0; i < tasksCount; i++ {
		// taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		r := rand.Intn(100)
		sumTime += Slow(r)

		tasks = append(tasks, func() error {
			// time.Sleep(taskSleep)
			r := r
			Slow(r)
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}

	start := time.Now()
	err := Run(tasks, workersCount, maxErrorsCount)
	elapsedTime := time.Since(start)
	require.NoError(t, err)

	require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
}

func withErrors(t *testing.T, tasksCount int, workersCount int, maxErrorsCount int) {
	t.Helper()
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			// time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			Slow(rand.Intn(100))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	err := Run(tasks, workersCount, maxErrorsCount)

	require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
}

func Slow(iterations int) time.Duration {
	const slowWeight = 1000000

	start := time.Now()
	a := 0
	for i := 0; i < iterations*slowWeight; i++ {
		a += i
	}
	return time.Since(start)
}
