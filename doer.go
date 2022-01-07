package doer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Processor struct {
	wg  *sync.WaitGroup
	sem *semaphore.Weighted
}

func New(maxWorkers int64) *Processor {
	return &Processor{
		wg:  &sync.WaitGroup{},
		sem: semaphore.NewWeighted(maxWorkers),
	}
}

// Do executes an action when it is allowed to.
func (un *Processor) Do(action func() error) {
	if err := un.sem.Acquire(context.Background(), 1); err != nil {
		panic(err)
	}
	un.wg.Add(1)

	go un.do(action)
}

//
func (un *Processor) do(action func() error) {
	defer un.wg.Done()
	defer un.sem.Release(1)

	err := action()
	if err != nil {
		fmt.Println(
			fmt.Sprintf(
				"error while processing: %w",
				err,
			),
		)
	}
}

func (un *Processor) Wait() error {
	un.wg.Wait()
	return nil
}
