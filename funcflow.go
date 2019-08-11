package funcrunner

import (
	"context"
)

// FuncFlow type
type FuncFlow interface {
	// ID returns identifier
	ID() string
	// Enabled?
	Enabled() (bool, string)
	// PreRun runs prior to Run(context.Context, []interface{}) error
	PreRun(context.Context) ([]interface{}, error)
	// Run is the main run func
	Run(context.Context, []interface{}) ([]interface{}, error)
	// PostRun runs after Run(context.Context, []interface{}) error
	PostRun(context.Context, []interface{}) error
}
