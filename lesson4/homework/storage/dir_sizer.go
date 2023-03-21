package storage

import (
	"context"
	"golang.org/x/sync/errgroup"
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
	//err []error
	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) worker(ctx context.Context, d Dir, result *Result) error {
	g, ctx := errgroup.WithContext(ctx)
	if dirList, fileList, err := d.Ls(ctx); err == nil {
		for _, i := range fileList {
			file := i
			g.Go(func() error {
				if delta, err1 := file.Stat(ctx); err1 == nil {
					atomic.AddInt64(&result.Size, delta)
					atomic.AddInt64(&result.Count, 1)
				} else {
					return err1
				}
				return nil
			})
		}
		for _, i := range dirList {
			dir := i
			g.Go(func() error {
				err := a.worker(ctx, dir, result)
				if err != nil {
					return err
				}
				return nil
			})
		}
		err3 := g.Wait()
		if err3 != nil {
			return err3
		}
	} else {
		return err
	}
	return nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := Result{}
	err := a.worker(ctx, d, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
