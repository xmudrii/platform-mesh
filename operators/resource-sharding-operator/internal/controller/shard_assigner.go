package controller

import (
	"sync"
	"sync/atomic"
)

type ShardAssigner struct {
	mu      sync.Mutex
	shards  []string
	counter atomic.Uint64
}

func NewShardAssigner(shards []string) *ShardAssigner {
	return &ShardAssigner{
		shards: shards,
	}
}

func (a *ShardAssigner) Next() string {
	idx := a.counter.Add(1) - 1
	return a.shards[idx%uint64(len(a.shards))]
}

func (a *ShardAssigner) UpdateShards(shards []string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shards = shards
}

func (a *ShardAssigner) RecoverFromDistribution(counts map[string]int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.shards) == 0 {
		return
	}

	minCount := -1
	minIdx := 0
	for i, shard := range a.shards {
		count, ok := counts[shard]
		if !ok {
			count = 0
		}
		if minCount < 0 || count < minCount {
			minCount = count
			minIdx = i
		}
	}

	a.counter.Store(uint64(minIdx))
}
