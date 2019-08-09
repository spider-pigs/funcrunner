# funcrunner
[![Build Status](https://travis-ci.org/spider-pigs/funcrunner.svg?branch=master)](https://travis-ci.org/spider-pigs/funcrunner) [![Go Report Card](https://goreportcard.com/badge/github.com/spider-pigs/funcrunner)](https://goreportcard.com/report/github.com/spider-pigs/funcrunner) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/69dc79e40fed443ba7b502de8a956bed)](https://www.codacy.com/app/spider-pigs/funcrunner?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=spider-pigs/funcrunner&amp;utm_campaign=Badge_Grade) [![GoDoc](https://godoc.org/github.com/spider-pigs/funcrunner?status.svg)](https://godoc.org/github.com/spider-pigs/funcrunner)

A small golang library that runs a set of functions safely (without panics) and returns simple stats.

## Install

```Go
import "github.com/spider-pigs/funcrunner"
```

## Usage

```Go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spider-pigs/funcrunner"
)

func main() {
	// create functions that return an error and takes context as param
	ctx := context.Background()
	f1 := func(ctx context.Context) error {
		return nil
	}
	f2 := func(ctx context.Context) error {
		return errors.New("I can't take it anymore")
	}

	// create the runner
	runner := funcrunner.Runner{}
	runner.FuncDone = func(f funcrunner.Func, duration time.Duration, err error) {
		fmt.Printf("function %s has completed, took %s.\n", f, duration)
	}

	// Start the runner
	duration, errored := runner.Run(ctx, f1, f2)
	fmt.Printf("runner has completed, took %s and %d/2 errored.\n", duration, errored)
}
```
