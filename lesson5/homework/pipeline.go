package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func doStage(ctx context.Context, in In) Out {
	out := make(chan any)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-in:
				if !ok {
					return
				}
				out <- n
			}
		}
	}()
	return out
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	empty := make(chan any)
	defer close(empty)
	for _, stage := range stages {
		if in == nil {
			return empty
		}
		in = stage(doStage(ctx, in))
	}
	return in
}
