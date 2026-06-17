package ratelimiter

import "sync"

type MemoryStore struct {
	states map[string]record
	mu     sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		states: make(map[string]record),
	}
}

func (s *MemoryStore) Get(key string) (memoryAlgorithmState, int, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, exists := s.states[key]

	if !exists {
		return nil, 0, false, nil
	}

	return rec.state, rec.version, true, nil
}

func (s *MemoryStore) CompareAndSwap(key string, version int, state memoryAlgorithmState) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, exists := s.states[key]

	if !exists {
		if version != 0 {
			return false, nil // Not an explicit error as this is a valid CAS conflict.
		}
		s.states[key] = record{state: state, version: version + 1}
		return true, nil
	}

	if version != rec.version {
		return false, nil
	}

	s.states[key] = record{state: state, version: version + 1}
	return true, nil
}
