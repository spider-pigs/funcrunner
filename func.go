package funcrunner

import (
	"context"
	"reflect"
	"runtime"
	"strings"
)

// Func type
type Func func(context.Context) error

func (f Func) String() string {
	p := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	split := strings.Split(p, ".")
	funcName := split[len(split)-1]
	return strings.TrimRight(funcName, "-fm")
}
