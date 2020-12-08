package funcrunner_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/spider-pigs/funcrunner"
)

var (
	happyFunc = func(ctx context.Context) error {
		return nil
	}
	sadFunc = func(ctx context.Context) error {
		return errors.New(":-(")
	}
)

type fflow struct{}

func (fflow fflow) ID() string                                                { return fmt.Sprintf("%T", fflow) }
func (fflow fflow) Enabled() (bool, string)                                   { return true, "" }
func (fflow fflow) PreRun(context.Context) ([]interface{}, error)             { return nil, nil }
func (fflow fflow) Run(context.Context, []interface{}) ([]interface{}, error) { return nil, nil }
func (fflow fflow) PostRun(context.Context, []interface{}) error              { return nil }
func (fflow fflow) String() string                                            { return fflow.ID() }

type fflowDisabled struct{}

func (fflow fflowDisabled) ID() string                                    { return fmt.Sprintf("%T", fflow) }
func (fflow fflowDisabled) Enabled() (bool, string)                       { return false, "some reason" }
func (fflow fflowDisabled) PreRun(context.Context) ([]interface{}, error) { return nil, nil }
func (fflow fflowDisabled) Run(context.Context, []interface{}) ([]interface{}, error) {
	return nil, nil
}
func (fflow fflowDisabled) PostRun(context.Context, []interface{}) error { return nil }
func (fflow fflowDisabled) String() string                               { return fflow.ID() }

func TestRunFuncs(t *testing.T) {
	fcount := 0
	runner := funcrunner.Runner{}
	runner.FuncDone = func(f funcrunner.Func, duration time.Duration, err error) {
		fcount++
	}
	ctx := context.Background()
	funcs := []funcrunner.Func{happyFunc, sadFunc}
	duration := runner.RunFuncs(ctx, funcs...)
	if duration == 0 {
		t.Error("duration should not be 0")
	}
	if fcount != len(funcs) {
		t.Errorf("expected FuncDone to be called %v, but was called %v times.", len(funcs), fcount)
	}
}

func TestRunFlows(t *testing.T) {
	fcount := 0
	runner := funcrunner.Runner{}
	runner.FlowDone = func(f funcrunner.FuncFlow, duration time.Duration, err error) {
		fcount++
	}
	ctx := context.Background()
	flows := []funcrunner.FuncFlow{fflow{}}
	duration := runner.RunFlows(ctx, flows...)
	if duration == 0 {
		t.Error("duration should not be 0")
	}
	if fcount != len(flows) {
		t.Errorf("expected FlowDone to be called %v, but was called %v times.", len(flows), fcount)
	}
}

func TestDisabledRunFlows(t *testing.T) {
	fcount := 0
	runner := funcrunner.Runner{}
	runner.FlowDone = func(f funcrunner.FuncFlow, duration time.Duration, err error) {
		fcount++
	}
	ctx := context.Background()
	flows := []funcrunner.FuncFlow{fflowDisabled{}}
	duration := runner.RunFlows(ctx, flows...)
	if duration != 0 {
		t.Errorf("duration should be 0 but was %v", duration)
	}
	if fcount != len(flows) {
		t.Errorf("expected FlowDone to be called %v, but was called %v times.", len(flows), fcount)
	}
}
