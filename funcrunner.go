package funcrunner

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Timeout for funcs
var Timeout = time.Second * 10

// FuncDone type
type FuncDone func(f Func, time time.Duration, err error)

// Runner type
type Runner struct {
	FuncDone FuncDone
}

// Run runs funcs.
func (runner Runner) Run(ctx context.Context, funcs ...Func) (time.Duration, int) {
	errored := 0
	var duration time.Duration

	for _, f := range funcs {
		elapsed, err := runFunc(ctx, f)
		if err != nil {
			errored++
		}
		duration += elapsed
		if runner.FuncDone != nil {
			runner.FuncDone(f, elapsed, err)
		}
	}

	return duration, errored
}

func runFunc(ctx context.Context, f Func) (time.Duration, error) {
	var err error
	var elapsed time.Duration

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				panicstr := fmt.Sprintf("%s", r)
				err = errors.New("test suite panic: " + panicstr)
			}
		}()
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()
		start := time.Now()
		err = f(ctx)
		elapsed = time.Since(start)
	}()
	wg.Wait()
	return elapsed, err
}
