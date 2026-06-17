package ratelimiter

type record struct {
	state   memoryAlgorithmState
	version int
}

type StateStore interface {
	Get(key string) (memoryAlgorithmState, int, bool, error)
	CompareAndSwap(key string, version int, state memoryAlgorithmState) (bool, error)
}
