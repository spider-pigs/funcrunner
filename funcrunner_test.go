package funcrunner_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/spider-pigs/funcrunner"
)

func TestFuncDone(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	runner := funcrunner.Runner{}
	runner.FuncDone = func(f funcrunner.Func, duration time.Duration, err error) {
		// Checks that this call is done
		wg.Done()
	}
	ctx := context.Background()
	f := func(ctx context.Context) error {
		return nil
	}
	runner.Run(ctx, f)

	wg.Wait()
}
