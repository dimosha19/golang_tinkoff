package storage

import (
	"context"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	//maxWorkersCount int
	err []error
	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) worker(ctx context.Context, d Dir, result *Result) {
	wg := sync.WaitGroup{}
	if dirList, fileList, err := d.Ls(ctx); err == nil {
		for _, i := range fileList {
			wg.Add(1)
			go func(i File) {
				defer wg.Done()
				if delta, err1 := i.Stat(ctx); err1 == nil {
					atomic.AddInt64(&result.Size, delta)
					atomic.AddInt64(&result.Count, 1)
				} else {
					a.err = append(a.err, err1)
					return
				}
			}(i)
		}
		for _, i := range dirList {
			wg.Add(1)
			go func(i Dir) {
				defer wg.Done()
				a.worker(ctx, i, result)
			}(i)
		}
		wg.Wait()
	} else {
		a.err = append(a.err, err)
		return
	}
	return
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := Result{}
	a.worker(ctx, d, &res)
	if len(a.err) == 0 {
		return res, nil
	}
	return res, a.err[0]
}
