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
	value  any
}

func NewPack(val any, num int) *Pack {
	return &Pack{num, 0, val}
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

func Collector(b []*Pack) <-chan any {
	res := make(chan any)
	go func() {
		defer close(res)
		sort.Slice(b, func(i, j int) bool {
			return b[i].seqNum <= b[j].seqNum
		})
		for i := range b {
			res <- b[i].value
		}
	}()
	return res
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	stp := make(chan any)
	wg := sync.WaitGroup{}
	type ans struct {
		MU   sync.Mutex
		data []*Pack
	}
	var c ans
	ticket := 0
	for i := range in {
		select {
		case <-ctx.Done():
			close(stp)
			return stp
		default:
			wg.Add(1)
			go func(i any, num int) {
				defer wg.Done()
				n := NewPack(i, num)
				n = n.Run(stages...)
				c.MU.Lock()
				c.data = append(c.data, n)
				c.MU.Unlock()
			}(i, ticket)
			ticket++
		}
	}
	wg.Wait()
	return Collector(c.data)
}
