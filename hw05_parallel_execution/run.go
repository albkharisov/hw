package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type routineData struct {
	f   chan Task
	err chan error
	wg  sync.WaitGroup
}

func routine(rd *routineData) {
	for {
		if f, ok := <-rd.f; ok {
			rd.err <- f()
		} else {
			break
		}
	}

	rd.wg.Done()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	resultErr := ErrErrorsLimitExceeded
	workersNumber := min(len(tasks), n)
	rd := routineData{
		f:   make(chan Task, workersNumber),
		err: make(chan error, workersNumber),
	}
	rd.wg.Add(workersNumber)

	tasksQueued := 0
	for ; tasksQueued < workersNumber; tasksQueued++ {
		go routine(&rd)
		rd.f <- tasks[tasksQueued]
	}

Exit:
	for {
		err := <-rd.err
		if m >= 0 && err != nil {
			m--
			if m == 0 {
				break Exit
			}
		}

		if tasksQueued < len(tasks) {
			rd.f <- tasks[tasksQueued]
			tasksQueued++
		} else {
			resultErr = nil
			break Exit
		}
	}

	close(rd.f)
	rd.wg.Wait()
	close(rd.err)

	return resultErr
}
