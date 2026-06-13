package ratelimiter

import "hash/fnv"

type ShardedMemoryStore struct {
	shards []*MemoryStore
}

const MAXSHARDSIZE = 1000

func NewShardedMemoryStore(shardCount int) (*ShardedMemoryStore, error) {
	if shardCount <= 0 || shardCount > MAXSHARDSIZE {
		return nil, ErrInvalidShardCount
	}

	shards := make([]*MemoryStore, shardCount)

	for i := range shardCount {
		shards[i] = NewMemoryStore()
	}

	return &ShardedMemoryStore{shards}, nil
}

func (s *ShardedMemoryStore) shardIndex(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % uint32(len(s.shards)))
}

func (s *ShardedMemoryStore) Get(key string) (algorithmState, int, bool, error) {
	idx := s.shardIndex(key)
	return s.shards[idx].Get(key)
}

func (s *ShardedMemoryStore) CompareAndSwap(key string, version int, state algorithmState) (bool, error) {
	idx := s.shardIndex(key)
	return s.shards[idx].CompareAndSwap(key, version, state)
}
