package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func doStage(ctx context.Context, in <-chan any) <-chan any {
	out := make(chan any)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return
		default:
			for n := range in {
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
		in = stage(doStage(ctx, in))
	}
	return in
}
