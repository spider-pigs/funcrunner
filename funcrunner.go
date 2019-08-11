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

// FlowDone type
type FlowDone func(f FuncFlow, time time.Duration, err error)

// Runner type
type Runner struct {
	FuncDone FuncDone
	FlowDone FlowDone
}

// RunFuncs runs funcs.
func (runner Runner) RunFuncs(ctx context.Context, funcs ...Func) time.Duration {
	var duration time.Duration

	for _, f := range funcs {
		elapsed, err := runFunc(ctx, f)
		duration += elapsed
		if runner.FuncDone != nil {
			runner.FuncDone(f, elapsed, err)
		}
	}

	return duration
}

// RunFlows runs flows.
func (runner Runner) RunFlows(ctx context.Context, runs ...FuncFlow) time.Duration {
	var duration time.Duration

	for _, r := range runs {
		var elapsed time.Duration
		var err error
		enabled, _ := r.Enabled()
		if enabled {
			elapsed, err = runFlow(ctx, r)
			if err != nil {
				elapsed = 0
			}
			duration += elapsed
		}
		if runner.FlowDone != nil {
			runner.FlowDone(r, elapsed, err)
		}
	}

	return duration
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
				err = errors.New("func panic: " + panicstr)
			}
		}()
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()
		start := time.Now()
		err = f(ctx)
		elapsed = time.Since(start)
		if err != nil {
			elapsed = 0
		}
	}()
	wg.Wait()
	return elapsed, err
}

func runFlow(ctx context.Context, r FuncFlow) (time.Duration, error) {
	var err error
	var duration time.Duration
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				panicstr := fmt.Sprintf("%s", r)
				err = errors.New("func panic: " + panicstr)
			}
		}()
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()

		// Run pre func
		var args []interface{}
		args, err = r.PreRun(ctx)
		if err != nil {
			return
		}

		// Run main func
		start := time.Now()
		args, err = r.Run(ctx, args)
		duration = time.Since(start)
		if err != nil {
			return
		}

		// Run post func
		err = r.PostRun(ctx, args)
	}()
	wg.Wait()
	return duration, err
}
