package executor

import (
	"context"
	"sort"
	"sync"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

type Pack struct {
	seqNum int
	stage  int
	ctx    context.Context
	value  any
}

func NewPack(ctx context.Context, val any, num int) *Pack {
	return &Pack{num, 0, ctx, val}
}

func (p *Pack) Push(stages ...Stage) *Pack {
	if p.stage < len(stages) {
		k := make(chan any, 1)
		k <- p.value
		close(k)
		t := <-stages[p.stage](k)
		p.value = t
		p.stage++
	} else {
		p.stage = -1
	}
	return p
}

func (p *Pack) Run(stages ...Stage) *Pack {
	for p.stage != -1 {
		p = p.Push(stages...)
	}
	return p
}

func Collector(c <-chan *Pack) chan any {
	res := make(chan any, 100)
	defer close(res)
	var b []*Pack
	for i := range c {
		b = append(b, i)
	}
	sort.Slice(b, func(i, j int) bool {
		return b[i].seqNum <= b[j].seqNum
	})
	for i := range b {
		res <- b[i].value
	}
	return res
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	stp := make(chan any)
	wg := sync.WaitGroup{}
	ticket := 0
	c := make(chan *Pack, 100)
	for i := range in {
		select {
		case <-ctx.Done():
			close(stp)
			return stp
		default:
			wg.Add(1)
			go func(i any, num int) {
				defer wg.Done()
				n := NewPack(ctx, i, num)
				c <- n.Run(stages...)
			}(i, ticket)
			ticket++
		}
	}
	wg.Wait()
	close(c)
	return Collector(c)
}
